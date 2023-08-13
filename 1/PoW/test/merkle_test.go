package test

import (
	trans "concepts/transaction"
	utils "concepts/utils"
	"strings"
	"testing"
	"time"
)

func Test_Merkle_Invalid(t *testing.T) {
	//? Create  transactions
	transactions := []*trans.Transaction{}
	_, err := Db.BuildMerkle(transactions)
	if !strings.Contains(err.Error(), "transactions cannot be empty") {
		t.Fatalf("it should contain this error message")
	}
	t.Log("Invalid Merkle test OK")
}

func Test_Merkle_Single(t *testing.T) {
	//? create single transaction (without cryptographic material)
	transaction := &trans.Transaction{}
	transaction.Id = utils.GenerateRandomString(50)
	transaction.Timestamp = time.Now().UTC().Unix()
	transactions := []*trans.Transaction{}
	transactions = append(transactions, transaction)
	//? Build the merkle tree
	root, err := Db.BuildMerkle(transactions)
	if err != nil {
		t.Fatalf("it should not throw a error at this point %v", err)
	}
	//? Verify database state
	if Db.Data[root][0].Id != transaction.Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.MerkleTree[root].Hash != root {
		t.Fatalf("the root of the merkletree hash should be the returned hash")
	}
	hash, err := transaction.HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this state %v", err)
	}
	if Db.MerkleTree[root].Left.Hash != hash {
		t.Fatalf("left hash should have the same value as the hash of the transaction")
	}
	if Db.MerkleTree[root].Right.Hash != hash {
		t.Fatalf("right hash should have the same value as the hash of the transaction")
	}
}

func Test_Merkle(t *testing.T) {
	//? Create the transactions
	transactions := []*trans.Transaction{}
	for i := 0; i < 5; i++ {
		newTransaction := trans.Transaction{}
		newTransaction.Id = utils.GenerateRandomString(50)
		newTransaction.Timestamp = time.Now().UTC().Unix()
		transactions = append(transactions, &newTransaction)
	}
	//? Create merkle
	root, err := Db.BuildMerkle(transactions)
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	//? Verify database state
	if Db.Data[root][0].Id != transactions[0].Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.Data[root][1].Id != transactions[1].Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.Data[root][2].Id != transactions[2].Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.Data[root][3].Id != transactions[3].Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.Data[root][4].Id != transactions[4].Id {
		t.Fatalf("the transaction should be on the db")
	}
	if Db.MerkleTree[root].Hash != root {
		t.Fatalf("the root of the merkletree hash should be the returned hash")
	}
	hash, err := transactions[0].HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this point")
	}
	if hash != Db.MerkleTree[root].Left.Left.Left.Hash {
		t.Fatalf("transacion 0 hash should be the same has the one on the merkle tree")
	}
	hash, err = transactions[1].HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this point")
	}
	if hash != Db.MerkleTree[root].Left.Left.Right.Hash {
		t.Fatalf("transacion 1 hash should be the same has the one on the merkle tree")
	}
	hash, err = transactions[2].HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this point")
	}
	if hash != Db.MerkleTree[root].Left.Right.Left.Hash {
		t.Fatalf("transacion 2 hash should be the same has the one on the merkle tree")
	}
	hash, err = transactions[3].HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this point")
	}
	if hash != Db.MerkleTree[root].Left.Right.Right.Hash {
		t.Fatalf("transacion 3 hash should be the same has the one on the merkle tree")
	}
	hash, err = transactions[4].HashTransaction()
	if err != nil {
		t.Fatalf("should not throw a error at this point")
	}
	if hash != Db.MerkleTree[root].Right.Left.Left.Hash {
		t.Fatalf("transacion 4 hash should be the same has the one on the merkle tree")
	}
	if hash != Db.MerkleTree[root].Right.Left.Right.Hash {
		t.Fatalf("transacion 4 hash should be the same has the one on the merkle tree")
	}
	if hash != Db.MerkleTree[root].Right.Right.Left.Hash {
		t.Fatalf("transacion 4 hash should be the same has the one on the merkle tree")
	}
	if hash != Db.MerkleTree[root].Right.Right.Right.Hash {
		t.Fatalf("transacion 4 hash should be the same has the one on the merkle tree")
	}
}
