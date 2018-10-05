package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Secret handling using barbican",
	Long: `Secret handling using OpenStack Barbican.
	Secrets are stored encrypted in local files, with the key being stored in
	Barbican. These files can be safely committed to version control.`,
}

var Debug bool
var Verbose bool
var Deployment string
var SecretsFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "D", false, "Output debug info")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Increase verbosity")
	RootCmd.PersistentFlags().StringVarP(&SecretsFile, "secret-file", "s", "secrets.yaml", "Secrets file to manage")

	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	if Debug {
		log.SetLevel(log.DebugLevel)
	} else if Verbose {
		log.SetLevel(log.InfoLevel)
	}
}
