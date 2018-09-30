package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// commitCmd represents the hello command
var commitCmd = &cobra.Command{
	Use:   "commit <deployment>",
	Short: "commit value for a given deployment",
	Long: `This command commits the given value using the key associated with
	the given project`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit called with %v", args)
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)

}
