package test

import (
	t "concepts/transaction"
	"testing"
	"time"
)

var GlobalTransaction *t.Transaction
var LastOwnerPrivateKey string = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg0+FM8TPxRxQQTW+Y\nzZIRXj0cB5RvI3TGZt10TBbr0JahRANCAAQJRt3rIUnmKuc+7iZNHw+oK8cJ1oCP\niNJB4mnu1Nx8/ixTXhbgLz68lE0FlP3m6xhyuEA77eojx3nzUD0zM1sk\n-----END PRIVATE KEY-----"

func TestMain(m *testing.M) {
	//? genesis transaction
	headTransaction := &t.Transaction{}
	headTransaction.Amount = 10
	headTransaction.LastOwnerPublicKey = "-----BEGIN CERTIFICATE-----\nMIICvjCCAmWgAwIBAgIUbKOW6YXYtvt1/XrCLyuk/jDoesQwCgYIKoZIzj0EAwIw\nczELMAkGA1UEBhMCUFQxDjAMBgNVBAgTBVBvcnRvMRAwDgYDVQQHEwdBbGlhZG9z\nMR4wHAYDVQQKExVVbml2ZXJzaWRhZGUgZG8gTWluaG8xDzANBgNVBAsTBmNsaWVu\ndDERMA8GA1UEAxMIaW50ZXJhZG0wHhcNMjMwNjA3MTcwMjAwWhcNMjQwODA5MTAy\nODAwWjBwMQswCQYDVQQGEwJQVDEOMAwGA1UECBMFUG9ydG8xEDAOBgNVBAcTB0Fs\naWFkb3MxHjAcBgNVBAoTFVVuaXZlcnNpZGFkZSBkbyBtaW5obzEPMA0GA1UECxMG\nY2xpZW50MQ4wDAYDVQQDEwVwZWRybzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBAlG3eshSeYq5z7uJk0fD6grxwnWgI+I0kHiae7U3Hz+LFNeFuAvPryUTQWU/ebr\nGHK4QDvt6iPHefNQPTMzWySjgdkwgdYwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB\n/wQCMAAwHQYDVR0OBBYEFJW/Bybo0PCHGcwisuYjmwXnu6ugMB8GA1UdIwQYMBaA\nFPezAlpnV1BfmrS98b8PBq/Abm/0MBwGA1UdEQQVMBOCEXBlZHJvLWdmNjN0aGlu\nOXNjMFgGCCoDBAUGBwgBBEx7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoiIiwi\naGYuRW5yb2xsbWVudElEIjoicGVkcm8iLCJoZi5UeXBlIjoiY2xpZW50In19MAoG\nCCqGSM49BAMCA0cAMEQCIHY5XvJywAlQmIWV0ivqXC4zDawQ/SOTIRRMTKS+1k+S\nAiBxd3dGJpYKPGM0ek6wwNHFHnZtZClwq/0NzJ99HbFu2w==\n-----END CERTIFICATE-----"
	headTransaction.LastOwnerSignature = nil
	headTransaction.NewOwnerPublicKey = ""
	headTransaction.Timestamp = time.Now().UTC().Unix()
	GlobalTransaction = headTransaction
	//? Putting the data of the first transaction
	m.Run()
}
