package main

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestFetchKey(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSecretKey(t)
	HandleGetSecretKey(t)
	HandleGetPayloadKey(t)

	key, nonce, err := fetchKey(client.ServiceClient(), "test")
	if err != nil {
		t.Fatalf("test failed : %v", err)
	}
	if key != "s928HkkJKVCGO1q1aIFq1iWG3ZDh6LB7utsZ1mRqjKg=" {
		t.Fatalf("got wrong key in payload : %v vs %v", key, GetPayloadResponse)
	}
	if nonce != "cWcmxHPcuG0O0hY3" {
		t.Fatalf("got wrong nonce in payload : %v vs %v", nonce, GetPayloadResponse)
	}
}

// GetResponse provides a Get result.
const GetResponse = `
{
    "algorithm": "aes",
    "bit_length": 256,
    "content_types": {
        "default": "text/plain"
    },
    "created": "2018-06-21T02:49:48",
    "creator_id": "5c70d99f4a8641c38f8084b32b5e5c0e",
    "expiration": null,
    "mode": "cbc",
    "name": "test",
    "secret_ref": "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
    "secret_type": "opaque",
    "status": "ACTIVE",
    "updated": "2018-06-21T02:49:48"
}`

// GetPayloadResponse provides a payload result.
const GetPayloadResponse = `s928HkkJKVCGO1q1aIFq1iWG3ZDh6LB7utsZ1mRqjKg=
cWcmxHPcuG0O0hY3`

const ListResponse = `
{
    "secrets": [
        {
            "algorithm": "aes",
            "bit_length": 256,
            "content_types": {
                "default": "text/plain"
            },
            "created": "2018-06-21T02:49:48",
            "creator_id": "5c70d99f4a8641c38f8084b32b5e5c0e",
            "expiration": null,
            "mode": "cbc",
            "name": "test",
            "secret_ref": "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
            "secret_type": "opaque",
            "status": "ACTIVE",
            "updated": "2018-06-21T02:49:48"
        }
    ],
    "total": 1
}`

func HandleListSecretKey(t *testing.T) {
	th.Mux.HandleFunc("/secrets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListResponse)
	})
}

func HandleGetSecretKey(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetResponse)
	})
}

func HandleGetPayloadKey(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c/payload", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetPayloadResponse)
	})
}
