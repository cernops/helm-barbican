package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

var update = flag.Bool("update", false, "update golden test data")

func TestEncrypt(t *testing.T) {

	tests, err := filepath.Glob("testdata/encrypt_*.yaml")
	if err != nil {
		t.Fatalf("failed to parse test data :: %v", err)
	}
	for _, test := range tests {
		keyfile := fmt.Sprintf("%v.key", test)
		encrypted := fmt.Sprintf("%v.enc", test)

		content, err := ioutil.ReadFile(test)
		if err != nil {
			t.Fatalf("failed to read file %v :: %v", test, err)
		}

		// update golden data for test if flag given
		if *update {
			key, nonce, err := newKey()
			fmt.Printf("%v\n%v\n\n", key, nonce)
			ioutil.WriteFile(keyfile, []byte(fmt.Sprintf("%v\n%v", key, nonce)), 0644)
			result, err := encrypt(key, nonce, content)
			if err != nil {
				t.Fatalf("failed to encrypt data %v :: %v", test, err)
			}
			ioutil.WriteFile(encrypted, result, 0644)
		}

		keycontent, err := ioutil.ReadFile(keyfile)
		if err != nil {
			t.Fatalf("failed to read key %v :: %v", test, err)
		}
		key := strings.Split(string(keycontent), "\n")
		result, err := encrypt(string(key[0]), string(key[1]), content)
		if err != nil {
			t.Fatalf("failed to encrypt data %v :: %v", test, err)
		}

		expected, err := ioutil.ReadFile(encrypted)
		if err != nil {
			t.Fatalf("failed to read encrypted data %v :: %v", test, err)
		}
		if bytes.Compare(result, expected) != 0 {
			t.Errorf("expected: %v :: result: %v", expected, result)
		}
	}

}

func TestDecrypt(t *testing.T) {

	tests, err := filepath.Glob("testdata/encrypt_*.yaml")
	if err != nil {
		t.Fatalf("failed to parse test data :: %v", err)
	}
	for _, test := range tests {
		keyfile := fmt.Sprintf("%v.key", test)
		encryptedfile := fmt.Sprintf("%v.enc", test)

		content, err := ioutil.ReadFile(encryptedfile)
		if err != nil {
			t.Fatalf("failed to read file %v :: %v", test, err)
		}
		keycontent, err := ioutil.ReadFile(keyfile)
		if err != nil {
			t.Fatalf("failed to read key %v :: %v", test, err)
		}
		key := strings.Split(string(keycontent), "\n")
		result, err := decrypt(string(key[0]), string(key[1]), string(content))
		if err != nil {
			t.Fatalf("failed to encrypt data %v :: %v", test, err)
		}

		expected, err := ioutil.ReadFile(test)
		if err != nil {
			t.Fatalf("failed to read encrypted data %v :: %v", test, err)
		}
		if bytes.Compare(result, expected) != 0 {
			t.Errorf("expected: %v :: result: %v", expected, result)
		}
	}

}

func Example_encrypt() {
	// Sample yaml content
	content := `
		group:
		  value: 1
		  other: 2
	`

	// Sample AES-256 GCM key and nonce, use output from newKey()
	key := "s928HkkJKVCGO1q1aIFq1iWG3ZDh6LB7utsZ1mRqjKg="
	nonce := "cWcmxHPcuG0O0hY3"

	result, err := encrypt(key, nonce, []byte(content))
	if err != nil {
		log.Fatalf("encrypt failed : %v", err)
	}
	fmt.Printf("%v", string(result))
	// Output:
	// EetswXTplHOz9LXemnt4cglcWhp9/Uv2vTif1kiSSzuWY/Gyp953iL1X7JMDTBtpAo5W0Bo=
}

func Example_decrypt() {
	// Content to decrypt must be base64 encoded
	b64content := "EetswXTplHOz9LXemnt4cglcWhp9/Uv2vTif1kiSSzuWY/Gyp953iL1X7JMDTBtpAo5W0Bo="

	// Sample AES-256 GCM key and nonce, use output from newKey()
	key := "s928HkkJKVCGO1q1aIFq1iWG3ZDh6LB7utsZ1mRqjKg="
	nonce := "cWcmxHPcuG0O0hY3"

	result, err := decrypt(key, nonce, b64content)
	if err != nil {
		log.Fatalf("decrypt failed : %v", err)
	}
	fmt.Printf("%v", string(result))
	// Output:
	// group:
	//		  value: 1
	//		  other: 2
}
