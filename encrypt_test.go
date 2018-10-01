package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
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
