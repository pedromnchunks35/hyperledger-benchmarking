package transaction

import (
	u "concepts/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// ? Simple representation of a transaction in a blockchain
type Transaction struct {
	//? Transaction id
	Id string `json:"id"`
	//? This is the receiver
	NewOwnerPublicKey string `json:"new_owner_public_key"`
	//? This is the sender public key
	LastOwnerPublicKey string `json:"last_owner_public_key"`
	//? This is the sender signature
	LastOwnerSignature *SignatureEcdsa `json:"last_owner_signature"`
	//? Amount of coins
	Amount float64 `json:"amount"`
	//? Timestamp
	Timestamp int64 `json:"timestamp"`
	//? Las transaction
	LastTransaction *Transaction `json:"last_transaction"`
}

// ? The Ecdsa represents the signature with 2 variables, r and s
type SignatureEcdsa struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

// ? Signature content
type Signature struct {
	//? The last transaction
	LastTransaction Transaction `json:"transaction"`
	//? New owner
	NewOwner string `json:"new_owner"`
}

// ? function to load a public key from string
func loadPublicKeyFromString(key string) (*ecdsa.PublicKey, error) {
	//? Retrieve block
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("error converting to block")
	}
	//? parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	//? get the public key and convert it to ecdsa public key
	edcsaPublicKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parse to edcsa algorithm not ok")
	}
	return edcsaPublicKey, nil
}

// ? function to load a private key from string
func loadPrivateKeyFromString(key string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("error to decode PEM block")
	}
	//? Convert it to private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateEcdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("cannot parse to ecdsa")
	}
	return privateEcdsaKey, nil
}

// ? Function to sign data
func signData(key *ecdsa.PrivateKey, data []byte) (*SignatureEcdsa, error) {
	//? Create the signature from the private key
	r, s, err := ecdsa.Sign(
		rand.Reader,
		key,
		data,
	)
	if err != nil {
		return nil, err
	}
	return &SignatureEcdsa{R: r, S: s}, nil
}

// ? Function that represents a new transaction
func (t *Transaction) NewTransaction(
	newOwnerPublicKey string,
	lastOwnerPublicKey string,
	lastOwnerPrivateKey string,
	amount float64,
) (*Transaction, error) {
	//? key
	key, err := loadPrivateKeyFromString(lastOwnerPrivateKey)
	if err != nil {
		return nil, err
	}
	//? Create the signature
	content := &Signature{}
	content.LastTransaction = *t
	content.NewOwner = newOwnerPublicKey
	encoded, err := json.Marshal(*content)
	if err != nil {
		return nil, err
	}
	signature, err := signData(key, encoded)
	if err != nil {
		return nil, err
	}
	//? Create new transaction
	newTransaction := &Transaction{}
	newTransaction.Amount = amount
	newTransaction.LastOwnerPublicKey = lastOwnerPublicKey
	newTransaction.NewOwnerPublicKey = newOwnerPublicKey
	newTransaction.Timestamp = time.Now().UTC().Unix()
	newTransaction.LastOwnerSignature = signature
	newTransaction.LastTransaction = t
	newTransaction.Id = u.GenerateRandomString(30)
	return newTransaction, nil
}

// ? Function to verify the content of the transaction using the public key
func (t Transaction) Verify() bool {
	//? hash once again the content of the last transaction along with the public key
	signature := &Signature{}
	signature.NewOwner = t.NewOwnerPublicKey
	signature.LastTransaction = *t.LastTransaction
	//? get public key
	ecdsaPublicKey, err := loadPublicKeyFromString(t.LastOwnerPublicKey)
	if err != nil {
		return false
	}
	//? Hash the content
	encoded, err := json.Marshal(*signature)
	if err != nil {
		return false
	}
	//? make the verification
	return ecdsa.Verify(
		ecdsaPublicKey,
		encoded,
		t.LastOwnerSignature.R,
		t.LastOwnerSignature.S,
	)
}

// ? Function to check if a public key matches a given private key
func CheckKeys(privateKey string, certificate string) (bool, error) {
	//? Public key retrieval
	ecdsaPublicKey, err := loadPublicKeyFromString(certificate)
	if err != nil {
		return false, err
	}
	//? private key retrieval
	key, err := loadPrivateKeyFromString(privateKey)
	if err != nil {
		return false, err
	}
	content := []byte("hello world")
	signature, err := signData(key, content)
	if err != nil {
		return false, err
	}
	return ecdsa.Verify(ecdsaPublicKey, content, signature.R, signature.S), nil
}

// ? Function to hash a transaction
func (t Transaction) HashTransaction() (string, error) {
	//? encode it
	encoded, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(encoded)
	//? return as string
	return string(hashed[:]), nil
}
