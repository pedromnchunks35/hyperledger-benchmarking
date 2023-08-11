package transaction

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"time"
)

// ? Simple representation of a transaction in a blockchain
type Transaction struct {
	//? This is the receiver
	NewOwnerPublicKey string
	//? This is the sender public key
	LastOwnerPublicKey string
	//? This is the sender signature
	LastOwnerSignature string
	//? Amount of coins
	Amount float64
	//? Timestamp
	Timestamp int64
}

// ? Signature content
type Signature struct {
	//? The last transaction
	lastTransaction Transaction
	//? New owner
	newOwner string
}

// ? function to load a private key from string
func loadPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	//? Convert it to private key
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ? Function to sign data
func signData(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	//? hash the data
	hashed := sha256.Sum256(data)
	//? Create the signature from the private key
	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		key,
		crypto.SHA256,
		hashed[:],
	)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// ? Function that represents a new transaction
func (t *Transaction) NewTransaction(
	newOwnerPublicKey string,
	lastOwnerPublicKey string,
	lastOwnerPrivateKey string,
	amount float64,
) (*Transaction, error) {
	// ? Last transaction
	lastTransaction := t
	//? key
	key, err := loadPrivateKeyFromString(lastOwnerPrivateKey)
	if err != nil {
		return nil, err
	}
	//? Create the signature
	content := &Signature{}
	content.lastTransaction = *lastTransaction
	content.newOwner = newOwnerPublicKey
	signature, err := signData(key, []byte(fmt.Sprintf("%v", content)))
	if err != nil {
		return nil, err
	}
	//? Create new transaction
	newTransaction := &Transaction{}
	newTransaction.Amount = amount
	newTransaction.LastOwnerPublicKey = lastOwnerPublicKey
	newTransaction.NewOwnerPublicKey = newOwnerPublicKey
	newTransaction.Timestamp = time.Now().UTC().Unix()
	newTransaction.LastOwnerSignature = string(signature)
	return newTransaction, nil
}
