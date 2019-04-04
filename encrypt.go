package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// encryptCmd represents the 'enc' command.
var encryptCmd = &cobra.Command{
	Use:   "enc [FILE]",
	Short: "encrypt secrets with barbican key",
	Long: `This command encrypts the contents of a given secrets yaml file.
	The resulting file can then be safely committed to version control.
	This is a low level command which most times is not required, with
	'view' and 'edit' being preferred.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretsFile := args[0]
		content, err := ioutil.ReadFile(secretsFile)
		if err != nil {
			log.Fatalf("encrypt failed : %v", err)
		}
		if b64Encoded(string(content)) {
			log.Fatal("content is empty or already encrypted")
		}

		client, err := newKeyManager()
		if err != nil {
			log.Fatalf("could not init client :: %v", err)
		}
		key, nonce, err := fetchKey(client, releaseName())
		if err != nil {
			log.Fatalf("could not fetch key : %v", err)
		}

		result, err := encrypt(key, nonce, content)
		err = ioutil.WriteFile(secretsFile, result, 0644)
		if err != nil {
			log.Fatalf("encrypt failed : %v", err)
		}

	},
}

// decryptCmd represents the 'dec' command
var decryptCmd = &cobra.Command{
	Use:   "dec [FILE]",
	Short: "decrypt secrets with barbican key",
	Long: `This command decrypts the contents of a given secrets yaml file.
	This is a low level command which most times is not required, with 'view'
	and 'edit' being preferred.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretsFile := args[0]
		content, err := ioutil.ReadFile(secretsFile)
		if err != nil {
			log.Fatalf("decrypt failed : %v", err)
		}
		if !b64Encoded(string(content)) {
			log.Fatal("not touching unencrypted content")
		}
		client, err := newKeyManager()
		if err != nil {
			log.Fatalf("could not init client :: %v", err)
		}
		key, nonce, err := fetchKey(client, releaseName())
		if err != nil {
			log.Fatalf("could not get key : %v", err)
		}
		plain, err := decrypt(key, nonce, string(content))
		if err != nil {
			log.Fatalf("decrypt failed : %v", err)
		}
		err = ioutil.WriteFile(secretsFile, plain, 0644)
		if err != nil {
			log.Fatalf("could not write file : %v", err)
		}
	},
}

// viewCmd represents the 'view' command.
var viewCmd = &cobra.Command{
	Use:   "view [FILE]",
	Short: "decrypt and display secrets",
	Long: `This command decrypts the contents of a given secrets yaml file,
	and displays them in stdout. The contents are never stored unencrypted.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretsFile := args[0]
		content, err := ioutil.ReadFile(secretsFile)
		if err != nil {
			log.Fatalf("decrypt failed : %v", err)
		}
		if b64Encoded(string(content)) {
			client, err := newKeyManager()
			if err != nil {
				log.Fatalf("could not init client :: %v", err)
			}
			key, nonce, err := fetchKey(client, releaseName())
			if err != nil {
				log.Fatalf("could not get key :: %v", err)
			}
			content, err = decrypt(key, nonce, string(content))
			if err != nil {
				log.Fatalf("decrypt failed : %v", err)
			}
		}
		fmt.Printf("%v", string(content))
	},
}

// editCmd
var editCmd = &cobra.Command{
	Use:   "edit [FILE]",
	Short: "edit secrets",
	Long: `This command launches the system configured editor with the
	contents of a given secrets yaml file. The contents are decrypted for
	editing and encrypted on exit.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretsFile := args[0]
		client, err := newKeyManager()
		if err != nil {
			log.Fatalf("could not init client :: %v", err)
		}
		key, nonce, err := fetchKey(client, releaseName())
		if err != nil {
			log.Fatalf("could not fetch key : %v", err)
		}
		content, err := ioutil.ReadFile(secretsFile)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("decrypt failed : %v", err)
		}
		if b64Encoded(string(content)) {
			content, err = decrypt(key, nonce, string(content))
			if err != nil {
				log.Fatalf("decrypt failed : %v", err)
			}
		}
		ed, err := NewEditor()
		if err != nil {
			log.Fatalf("failed to find editor %v", err)
		}
		result, _, err := ed.LaunchTemp(strings.NewReader(string(content)))
		if err != nil {
			log.Fatalf("failed to open tmp file : %v", err)
		}
		encrypted, err := encrypt(key, nonce, result)
		if err != nil {
			log.Fatalf("failed to encrypt contents : %v", err)
		}
		err = ioutil.WriteFile(secretsFile, encrypted, 0600)
		if err != nil {
			log.Fatalf("failed to encrypt : %v", err)
		}
	},
}

func newKey() (string, string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(nonce), nil
}

func encrypt(b64key string, b64nonce string, payload []byte) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(b64nonce)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	sealed := aesgcm.Seal(nil, nonce, payload, nil)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(sealed)))
	base64.StdEncoding.Encode(result, sealed)
	return result, nil
}

func decrypt(b64key string, b64nonce string, b64payload string) ([]byte, error) {
	if b64payload == "" {
		return []byte{}, nil
	}
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(b64nonce)
	if err != nil {
		return nil, err
	}
	payload, err := base64.StdEncoding.DecodeString(b64payload)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plain, err := aesgcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func b64Encoded(content string) bool {
	_, err := base64.StdEncoding.DecodeString(content)
	if err == nil {
		return true
	}
	return false
}

func releaseName() string {
	if Release != "" {
		return Release
	}
	d, _ := os.Getwd()
	return filepath.Base(d)
}

func init() {
	RootCmd.AddCommand(encryptCmd)
	RootCmd.AddCommand(decryptCmd)
	RootCmd.AddCommand(viewCmd)
	RootCmd.AddCommand(editCmd)
}
