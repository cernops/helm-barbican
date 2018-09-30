package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// encryptCmd represents the hello command
var encryptCmd = &cobra.Command{
	Use:   "encrypt <deployment> <value>",
	Short: "encrypt value for a given deployment",
	Long: `This command encrypts the given value using the key associated with
	the given deployment`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("encrypt called with %v", args)
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)

}
