package test

import (
	block "concepts/block"
	trans "concepts/transaction"
	"strings"
	"testing"
)

func Test_Invalid_New_Block(t *testing.T) {
	newBlock := &block.Block{}
	newBlock.BlockNumber = 1
	tran := make([]*trans.Transaction, 11)
	for i := range tran {
		tran[i] = &trans.Transaction{}
	}
	_, err := newBlock.NewBlock(6, tran)
	if !strings.Contains(err.Error(), "the limit is 10 transactions") {
		t.Fatalf("it needs to throw a error")
	}
	t.Log("block invalid test OK")
}

func Test_New_Block(t *testing.T) {
	//? Genesis block
	criteria := 3
	genesis := &block.Block{}
	genesis.BlockNumber = 1
	//? append transaction
	tran := []*trans.Transaction{}
	tran = append(tran, GlobalTransaction)
	//? create new block
	result, err := genesis.NewBlock(criteria, tran)
	//? make verification
	if err != nil {
		t.Fatalf("it should not return a error, the block is valid")
	}
	newBlock := result
	if newBlock.BlockNumber != 2 {
		t.Fatalf("the block number must be 2 now")
	}
	if len(newBlock.Nonce) == 0 {
		t.Fatalf("the nonce should be a new random value")
	}
	hash, err := newBlock.HashBlock()
	if err != nil {
		t.Fatalf("it should not throw a error when hashing the block")
	}
	if !block.VerifyBlockHash([]byte(hash), criteria) {
		t.Fatalf("it should be a valid nonce")
	}
	if newBlock.PreviousBlock != genesis {
		t.Fatalf("previous block must be genesis")
	}
	if newBlock.Timestamp == 0 {
		t.Fatalf("timestamp should be a valid value")
	}
	if newBlock.Transactions[0] != tran[0] {
		t.Fatalf("it should have the Global transaction inside of it")
	}
	t.Log("block test OK")
}
