package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// decryptCmd represents the hello command
var decryptCmd = &cobra.Command{
	Use:   "decrypt <deployment> <value>",
	Short: "decrypt value for a given deployment",
	Long: `This command decrypts the given value using the key associated with
	the given deployment`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("decrypt called with %v", args)
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)

}
