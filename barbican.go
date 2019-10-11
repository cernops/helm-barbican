package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

func newKeyManager() (*gophercloud.ServiceClient, error) {
	opts, err := clientconfig.AuthOptions(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate :: %v : %v", opts, err)
	}
	provider, err := openstack.AuthenticatedClient(*opts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate :: %v : %v", opts, err)
	}

	client, err := openstack.NewKeyManagerV1(provider,
		gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	if err != nil {
		return nil, fmt.Errorf("failed to create key manager :: %v", err)
	}
	return client, nil
}

func fetchKey(client *gophercloud.ServiceClient, deployment string) (string, string, error) {
	// check secret exists, create if not
	pages, err := secrets.List(client, secrets.ListOpts{Name: deployment}).AllPages()
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
		secret, err := secrets.Create(client, createOpts).Extract()
		if err != nil {
			return "", "", err
		}
		secs = []secrets.Secret{*secret}
	}
	secretID, err := parseID(secs[0].SecretRef)
	if err != nil {
		return "", "", err
	}
	payload, err = secrets.GetPayload(client, secretID, nil).Extract()
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
