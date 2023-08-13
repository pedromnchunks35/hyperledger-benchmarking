package test

import (
	"bytes"
	trans "concepts/transaction"
	"encoding/json"
	"strings"
	"testing"
)

// ? Create new transaction
var newOwnerPublicKey string = "-----BEGIN CERTIFICATE-----\nMIICwzCCAmmgAwIBAgIUWvSEIV9mvdzV2KegUEhTPnYCwrwwCgYIKoZIzj0EAwIw\nczELMAkGA1UEBhMCUFQxDjAMBgNVBAgTBVBvcnRvMRAwDgYDVQQHEwdBbGlhZG9z\nMR4wHAYDVQQKExVVbml2ZXJzaWRhZGUgZG8gTWluaG8xDzANBgNVBAsTBmNsaWVu\ndDERMA8GA1UEAxMIaW50ZXJhZG0wHhcNMjMwNjA3MTcwMjAwWhcNMjQwNjE0MTM0\nMzAwWjByMQswCQYDVQQGEwJQVDEOMAwGA1UECBMFUG9ydG8xEDAOBgNVBAcTB0Fs\naWFkb3MxHjAcBgNVBAoTFVVuaXZlcnNpZGFkZSBkbyBtaW5obzEOMAwGA1UECxMF\nYWRtaW4xETAPBgNVBAMTCGFkbS1pdGVyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcD\nQgAEQYCC8W7yYGeoW7sBMJowdTlBbiO2VGCnRK+LBpZXiFw5OHbRAa563aKtqbiE\nop0Sqq7txjFCN8x4amXxmcfW8qOB2zCB2DAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0T\nAQH/BAIwADAdBgNVHQ4EFgQUSggva3G5J8NU5n4SKoXdRjL0VRcwHwYDVR0jBBgw\nFoAU97MCWmdXUF+atL3xvw8Gr8Bub/QwHAYDVR0RBBUwE4IRcGVkcm8tZ2Y2M3Ro\naW45c2MwWgYIKgMEBQYHCAEETnsiYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiIi\nLCJoZi5FbnJvbGxtZW50SUQiOiJhZG0taXRlciIsImhmLlR5cGUiOiJhZG1pbiJ9\nfTAKBggqhkjOPQQDAgNIADBFAiEA9zPQSc72DS1Gt7KYhkhJKZGid6Nwjk+c37ey\ntYPTgy8CID0tCHuV2FUxI6F+3eGhf9/MFepC8jiXA+87BQljE1ca\n-----END CERTIFICATE-----"

func Test_keys(t *testing.T) {
	ok, err := trans.CheckKeys(LastOwnerPrivateKey, GlobalTransaction.LastOwnerPublicKey)
	if err != nil {
		t.Fatalf("should not throw error %v", err)
	}
	if !ok {
		t.Fatalf("the keys must be a valid pair")
	}
	t.Log("valid keys test OK")
}

func Test_AddingTransaction(t *testing.T) {
	GlobalTransaction, err := GlobalTransaction.NewTransaction(
		newOwnerPublicKey,
		GlobalTransaction.LastOwnerPublicKey,
		LastOwnerPrivateKey,
		20,
	)
	if err != nil {
		t.Fatalf("it should not return a error at this stage %v", err)
	}
	//? Verification of the new transaction
	if !GlobalTransaction.Verify() {
		t.Fatalf("transaction should be authentic")
	}
	//? Get the content of the last transaction and display it
	var formattedJson *bytes.Buffer = bytes.NewBuffer([]byte{})
	encodedJson, err := json.Marshal(GlobalTransaction.LastTransaction)
	if err != nil {
		t.Fatalf("should not throw a error %v", err)
	}
	err = json.Indent(formattedJson, encodedJson, "", " ")
	if err != nil {
		t.Fatalf("should not throw a error %v", err)
	}
	//? Check last transaction data
	if GlobalTransaction.LastTransaction.Amount != 10 {
		t.Fatalf("it should be 10 the amount of transfer")
	}
	if !strings.Contains(GlobalTransaction.LastTransaction.LastOwnerPublicKey, "-----BEGIN CERTIFICATE-----\nMIICvjCCAmWgAwIBAgIUbKOW6YXYtvt1/XrCLyuk/jDoesQwCgYIKoZIzj0EAwIw\nczELMAkGA1UEBhMCUFQxDjAMBgNVBAgTBVBvcnRvMRAwDgYDVQQHEwdBbGlhZG9z\nMR4wHAYDVQQKExVVbml2ZXJzaWRhZGUgZG8gTWluaG8xDzANBgNVBAsTBmNsaWVu\ndDERMA8GA1UEAxMIaW50ZXJhZG0wHhcNMjMwNjA3MTcwMjAwWhcNMjQwODA5MTAy\nODAwWjBwMQswCQYDVQQGEwJQVDEOMAwGA1UECBMFUG9ydG8xEDAOBgNVBAcTB0Fs\naWFkb3MxHjAcBgNVBAoTFVVuaXZlcnNpZGFkZSBkbyBtaW5obzEPMA0GA1UECxMG\nY2xpZW50MQ4wDAYDVQQDEwVwZWRybzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBAlG3eshSeYq5z7uJk0fD6grxwnWgI+I0kHiae7U3Hz+LFNeFuAvPryUTQWU/ebr\nGHK4QDvt6iPHefNQPTMzWySjgdkwgdYwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB\n/wQCMAAwHQYDVR0OBBYEFJW/Bybo0PCHGcwisuYjmwXnu6ugMB8GA1UdIwQYMBaA\nFPezAlpnV1BfmrS98b8PBq/Abm/0MBwGA1UdEQQVMBOCEXBlZHJvLWdmNjN0aGlu\nOXNjMFgGCCoDBAUGBwgBBEx7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoiIiwi\naGYuRW5yb2xsbWVudElEIjoicGVkcm8iLCJoZi5UeXBlIjoiY2xpZW50In19MAoG\nCCqGSM49BAMCA0cAMEQCIHY5XvJywAlQmIWV0ivqXC4zDawQ/SOTIRRMTKS+1k+S\nAiBxd3dGJpYKPGM0ek6wwNHFHnZtZClwq/0NzJ99HbFu2w==\n-----END CERTIFICATE-----") {
		t.Fatalf("it should contain the correct public key")
	}
	if GlobalTransaction.LastTransaction.LastOwnerSignature != nil {
		t.Fatalf("it must be a genesis transaction")
	}
	if !strings.Contains(GlobalTransaction.LastTransaction.NewOwnerPublicKey, "") {
		t.Fatalf("genesis transaction should have newowner as empty")
	}
	if GlobalTransaction.LastTransaction.Timestamp == 0 {
		t.Fatalf("it should give a correct timestamp")
	}
	t.Logf("%v\n", formattedJson.String())
	//? compare common things
	if !strings.Contains(GlobalTransaction.LastTransaction.LastOwnerPublicKey, GlobalTransaction.LastOwnerPublicKey) {
		t.Fatalf("they should be the same")
	}
	//? Get the content of current transaction and display it
	encodedJson, err = json.Marshal(GlobalTransaction)
	if err != nil {
		t.Fatalf("should not throw a error %v", err)
	}
	err = json.Indent(formattedJson, encodedJson, "", " ")
	if err != nil {
		t.Fatalf("should not throw a error %v", err)
	}
	//? Check current transaction data
	if GlobalTransaction.Amount != 20 {
		t.Fatalf("it should be 20 the amount of transfer")
	}
	if !strings.Contains(GlobalTransaction.LastOwnerPublicKey, "-----BEGIN CERTIFICATE-----\nMIICvjCCAmWgAwIBAgIUbKOW6YXYtvt1/XrCLyuk/jDoesQwCgYIKoZIzj0EAwIw\nczELMAkGA1UEBhMCUFQxDjAMBgNVBAgTBVBvcnRvMRAwDgYDVQQHEwdBbGlhZG9z\nMR4wHAYDVQQKExVVbml2ZXJzaWRhZGUgZG8gTWluaG8xDzANBgNVBAsTBmNsaWVu\ndDERMA8GA1UEAxMIaW50ZXJhZG0wHhcNMjMwNjA3MTcwMjAwWhcNMjQwODA5MTAy\nODAwWjBwMQswCQYDVQQGEwJQVDEOMAwGA1UECBMFUG9ydG8xEDAOBgNVBAcTB0Fs\naWFkb3MxHjAcBgNVBAoTFVVuaXZlcnNpZGFkZSBkbyBtaW5obzEPMA0GA1UECxMG\nY2xpZW50MQ4wDAYDVQQDEwVwZWRybzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBAlG3eshSeYq5z7uJk0fD6grxwnWgI+I0kHiae7U3Hz+LFNeFuAvPryUTQWU/ebr\nGHK4QDvt6iPHefNQPTMzWySjgdkwgdYwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB\n/wQCMAAwHQYDVR0OBBYEFJW/Bybo0PCHGcwisuYjmwXnu6ugMB8GA1UdIwQYMBaA\nFPezAlpnV1BfmrS98b8PBq/Abm/0MBwGA1UdEQQVMBOCEXBlZHJvLWdmNjN0aGlu\nOXNjMFgGCCoDBAUGBwgBBEx7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoiIiwi\naGYuRW5yb2xsbWVudElEIjoicGVkcm8iLCJoZi5UeXBlIjoiY2xpZW50In19MAoG\nCCqGSM49BAMCA0cAMEQCIHY5XvJywAlQmIWV0ivqXC4zDawQ/SOTIRRMTKS+1k+S\nAiBxd3dGJpYKPGM0ek6wwNHFHnZtZClwq/0NzJ99HbFu2w==\n-----END CERTIFICATE-----") {
		t.Fatalf("it should contain the correct public key")
	}
	if GlobalTransaction.LastOwnerSignature == nil {
		t.Fatalf("it must have the correct lastOwner Signature, this isnt genesis")
	}
	if !strings.Contains(GlobalTransaction.NewOwnerPublicKey, newOwnerPublicKey) {
		t.Fatalf("It should have a new owner")
	}
	if GlobalTransaction.Timestamp == 0 {
		t.Fatalf("it should give a correct timestamp")
	}
	t.Logf("%v\n", formattedJson.String())
	t.Logf("transaction test OK")
}
