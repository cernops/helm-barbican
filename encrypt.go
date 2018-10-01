package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/spf13/cobra"
)

var Deployment string
var SecretsFile string

// encryptCmd represents the hello command
var encryptCmd = &cobra.Command{
	Use:   "enc",
	Short: "encrypt value for a given deployment",
	Long: `This command encrypts the given value using the key associated with
	the given deployment`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := ioutil.ReadFile(SecretsFile)
		if err != nil {
			fmt.Printf("encrypt failed %v\n", err)
		}

		key, nonce, err := deploymentKey(Deployment)
		if err != nil {
			fmt.Printf("failed to get deployment key :: %v :: %v\n", Deployment, err)
		}

		result, err := encrypt(key, nonce, content)
		err = ioutil.WriteFile(SecretsFile, result, 0644)
		if err != nil {
			fmt.Printf("encrypt failed :: %v\n", err)
		}

	},
}

// decryptCmd represents the hello command
var decryptCmd = &cobra.Command{
	Use:   "dec",
	Short: "decrypt value for a given deployment",
	Long: `This command decrypts the given value using the key associated with
	the given deployment`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := ioutil.ReadFile(SecretsFile)
		if err != nil {
			fmt.Printf("decrypt failed :: %v\n", err)
		}
		key := "9OhJKoy5Y2urelQrD7tbyy5FJi7nBbbQvvvmuaqJSj0="
		nonce := "aBVOqBxSc++tWIa1"
		plain, err := decrypt(key, nonce, string(content))
		if err != nil {
			fmt.Printf("decrypt failed :: %v\n", err)
		}
		err = ioutil.WriteFile(SecretsFile, plain, 0644)
		if err != nil {
			fmt.Printf("decrypt failed :: %v\n", err)
		}
	},
}

// viewCmd
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "decrypt value for a given deployment",
	Long: `This command decrypts the given value using the key associated with
	the given deployment`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := ioutil.ReadFile(SecretsFile)
		if err != nil {
			fmt.Printf("decrypt failed :: %v\n", err)
		}
		key := "9OhJKoy5Y2urelQrD7tbyy5FJi7nBbbQvvvmuaqJSj0="
		nonce := "aBVOqBxSc++tWIa1"
		plain, err := decrypt(key, nonce, string(content))
		if err != nil {
			fmt.Printf("decrypt failed :: %v\n", err)
		}
		fmt.Printf("%v\n", string(plain))
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

func newKeyManager() (*gophercloud.ServiceClient, error) {
	opts, err := openstack.AuthOptionsFromEnv()
	opts.DomainID = "default"
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate :: %v", err)
	}

	client, err := openstack.NewKeyManagerV1(provider,
		gophercloud.EndpointOpts{Region: "cern"})
	if err != nil {
		return nil, fmt.Errorf("failed to create key manager :: %v", err)
	}
	return client, nil
}

func deploymentKey(deployment string) (string, string, error) {
	km, err := newKeyManager()
	if err != nil {
		return "", "", err
	}
	// check secret exists, create if not
	pages, err := secrets.List(km, secrets.ListOpts{Name: deployment}).AllPages()
	if err != nil {
		return "", "", err
	}
	secs, err := secrets.ExtractSecrets(pages)
	if err != nil {
		return "", "", err
	}
	var payload []byte
	if len(secs) == 0 {
		key, nonce, err := newKey()
		if err != nil {
			return "", "", err
		}
		payload = []byte(fmt.Sprintf("%v\n%v", key, nonce))
		createOpts := secrets.CreateOpts{
			Algorithm:          "aes",
			BitLength:          256,
			Mode:               "gcm",
			Name:               deployment,
			Payload:            string(payload),
			PayloadContentType: "text/plain",
			SecretType:         secrets.OpaqueSecret,
		}
		secret, err := secrets.Create(km, createOpts).Extract()
		if err != nil {
			return "", "", err
		}
		secs = []secrets.Secret{*secret}
	}
	secretID, err := parseID(secs[0].SecretRef)
	if err != nil {
		return "", "", err
	}
	payload, err = secrets.GetPayload(km, secretID, nil).Extract()
	if err != nil {
		return "", "", err
	}
	key := strings.Split(string(payload), "\n")
	return key[0], key[1], nil
}

func parseID(ref string) (string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("Could not parse %s", ref)
	}

	return parts[len(parts)-1], nil
}

func init() {
	encryptCmd.Flags().StringVarP(&Deployment, "deployment", "", "", "Destination deployment for this value")
	encryptCmd.Flags().StringVarP(&SecretsFile, "secret-file", "s", "secrets.yaml", "Secrets file to encrypt/decrypt")
	RootCmd.AddCommand(encryptCmd)
	RootCmd.AddCommand(decryptCmd)
	RootCmd.AddCommand(viewCmd)
}
