package main

import (
	"os"
	"strings"

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
var Release string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Cleanup OpenStack environment for compatibility with the golang client
	for _, v := range []string{
		"OS_IDENTITY_PROVIDER", "OS_AUTH_TYPE", "OS_MUTUAL_AUTH", "OS_PROTOCOL"} {
		os.Unsetenv(v)
	}
	os.Setenv("OS_AUTH_URL", strings.Replace(os.Getenv("OS_AUTH_URL"), "krb/", "", 1))
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "enable verbose output")
	RootCmd.PersistentFlags().StringVarP(&Release, "name", "n", "", "release name - if unspecified, the current directory name")

	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	if Debug {
		log.SetLevel(log.DebugLevel)
	} else if Verbose {
		log.SetLevel(log.InfoLevel)
	}
}
