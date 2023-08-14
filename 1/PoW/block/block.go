package block

import (
	merkle "concepts/merkle"
	trans "concepts/transaction"
	u "concepts/utils"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type BlockV2 struct {
	Nonce         string `json:"nonce"`
	BlockNumber   int    `json:"block_number"`
	Merkle        string `json:"merkle_hash"`
	PreviousBlock *Block `json:"previous_block"`
	Timestamp     int64  `json:"timestamp"`
}

// ? Function to verify a given block
func (b BlockV2) VerifyHashInBlock(transactionHash string, db *merkle.Database) (bool, error) {
	merkle := db.MerkleTree[b.Merkle]
	return merkle.VerifyHash(transactionHash)
}

type Block struct {
	Nonce         string               `json:"nonce"`
	BlockNumber   int                  `json:"block_number"`
	Transactions  []*trans.Transaction `json:"transactions"`
	PreviousBlock *Block               `json:"previous_block"`
	Timestamp     int64                `json:"timestamp"`
}

// ? Function to hash the block
func (b Block) HashBlock() (string, error) {
	//? Encode the block
	encoded, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(encoded)
	return string(hashed[:]), nil
}

// ? Function to verify the hash
func VerifyBlockHash(newHash []byte, criteria int) bool {
	//? Loop all over the hash until find 5 zeros, case it does not find them, then we need to generate another nonce
	count := 0
	for i := range newHash {
		if count == criteria {
			return true
		}
		if newHash[i] == byte(0) {
			count++
		}
	}
	return false
}

// ? Function to create a new block, given a number of zeros (criteria)
func (b *Block) NewBlock(criteria int, transactions []*trans.Transaction) (*Block, error) {
	if len(transactions) > 10 {
		return nil, fmt.Errorf("the limit is 10 transactions")
	}
	//? Create a new block
	newBlock := &Block{}
	newBlock.BlockNumber = b.BlockNumber + 1
	newBlock.PreviousBlock = b
	newBlock.Timestamp = time.Now().UTC().Unix()
	newBlock.Transactions = transactions
	for {
		//? Get the nonce
		nonce := u.GenerateRandomString(50)
		newBlock.Nonce = nonce
		//? Hash the content
		hash, err := newBlock.HashBlock()
		if err != nil {
			return nil, err
		}
		//? verify the hash
		if VerifyBlockHash([]byte(hash), criteria) {
			break
		}
	}
	b = newBlock
	return b, nil
}
