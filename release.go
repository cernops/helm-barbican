package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type Chart struct {
	Cluster string `yaml:"cluster"`
	Project string `yaml:"project"`
}

var chart *Chart

// configCmd
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set the current deployment release",
	Long:  `This command `,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		chart, err = chartYaml()
		if err != nil {
			log.Fatalf("failed to parse chart config : %v", err)
		}
		content, err := ioutil.ReadFile(SecretsFile)
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
		key, nonce, err := fetchKey(client, Deployment)
		if err != nil {
			log.Fatalf("could not fetch key : %v", err)
		}

		result, err := encrypt(key, nonce, content)
		err = ioutil.WriteFile(SecretsFile, result, 0644)
		if err != nil {
			log.Fatalf("encrypt failed : %v", err)
		}

	},
}

func chartYaml() (*Chart, error) {

	var chart Chart

	yamlFile, err := ioutil.ReadFile("Chart.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &chart)
	if err != nil {
		return nil, err
	}

	return &chart, nil
}

func init() {
	RootCmd.AddCommand(setCmd)
}
