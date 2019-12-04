package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// installCmd wraps the kubectl 'apply' command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "wrapper for kubectl apply, decrypting secrets",
	Long: `This command wraps the default kubectl apply command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapKubectlCommand("apply", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

// installCmd wraps the kubectl 'create' command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "wrapper for kubectl create, decrypting secrets",
	Long: `This command wraps the default kubectl create command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapKubectlCommand("create", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

// installCmd wraps the helm 'install' command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "wrapper for helm install, decrypting secrets",
	Long: `This command wraps the default helm install command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapHelmCommand("install", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

// upgradeCmd wraps the helm 'upgrade' command.
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "wrapper for helm upgrade, decrypting secrets",
	Long: `This command wraps the default helm upgrade command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapHelmCommand("upgrade", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

// lintCmd wraps the helm 'lint' command.
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "wrapper for helm lint, decrypting secrets",
	Long: `This command wraps the default helm lint command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapHelmCommand("lint", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

// templateCmd wraps the helm 'template' command.
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "wrapper for helm template, decrypting secrets",
	Long: `This command wraps the default helm template command,
	but decrypting any encrypted values file using Barbican. Available
	arguments are the same as for the default command.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := wrapHelmCommand("template", args)
		if err != nil {
			log.Fatalf("%v", string(out))
		}
		fmt.Printf(string(out))
	},
}

func wrapHelmCommand(cmd string, args []string) ([]byte, error) {
	for i, flag := range args {
		if i == len(args)-1 {
			// Last arg, can't possibly be --name and value remaining
			break
		}

		if flag == "--name" || flag == "-n" {
			Release = args[i+1]
		}
	}
	helmArgs, decryptedFiles, err := decryptSecrets(args)
	for _, f := range decryptedFiles {
		defer os.Remove(f)
	}
	if err != nil {
		return []byte{}, err
	}
	fullArgs := append([]string{cmd}, helmArgs...)
	helmCmd := exec.Command("helm", fullArgs...)
	return helmCmd.CombinedOutput()
}

func wrapKubectlCommand(cmd string, args []string) ([]byte, error) {
	helmArgs, decryptedFiles, err := decryptSecrets(args)
	for _, f := range decryptedFiles {
		defer os.Remove(f)
	}
	if err != nil {
		return []byte{}, err
	}
	fullArgs := append([]string{cmd}, helmArgs...)
	helmCmd := exec.Command("kubectl", fullArgs...)
	return helmCmd.CombinedOutput()
}

func decryptSecrets(args []string) ([]string, []string, error) {
	decryptedFiles := []string{}
	helmArgs := args
	for i, flag := range args {
		if flag == "--values" || flag == "-f" || flag == "--filename" {
			if len(helmArgs) > i+1 {
				fname := helmArgs[i+1]
				// Move to next arg if it does not exist
				content, err := ioutil.ReadFile(fname)
				if _, err := os.Stat(fname); os.IsNotExist(err) {
					continue
				}
				// Check if content is b64encoded, if not move on
				if !b64Encoded(string(content)) {
					continue
				}
				// Decrypt the contents
				client, err := newKeyManager()
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				key, nonce, err := fetchKey(client, releaseName())
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				plain, err := decrypt(key, nonce, string(content))
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				// Store decrypted contents in a shm file
				uuid, err := uuid.NewRandom()
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				tmpf := fmt.Sprintf("/dev/shm/%v", uuid)
				decryptedFiles = append(decryptedFiles, tmpf)
				_, err = os.OpenFile(tmpf, os.O_RDWR|os.O_CREATE, 0600)
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				err = ioutil.WriteFile(tmpf, plain, 0644)
				if err != nil {
					return helmArgs, decryptedFiles, err
				}
				// Update args to access the decrypt shm file instead
				helmArgs[i+1] = tmpf
			}
		}
	}
	return helmArgs, decryptedFiles, nil
}

func init() {
	if strings.Contains(os.Args[0], "kubectl-") {
		RootCmd.AddCommand(applyCmd)
		RootCmd.AddCommand(createCmd)
	} else {
		RootCmd.AddCommand(installCmd)
		RootCmd.AddCommand(upgradeCmd)
		RootCmd.AddCommand(lintCmd)
		RootCmd.AddCommand(templateCmd)
	}
}
