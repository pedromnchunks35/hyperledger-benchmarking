package merkle

import (
	t "concepts/transaction"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
)

// ? Database of Merkle Tres per group of transactions
type Database struct {
	MerkleTree map[string]*Merkle          `json:"merkle_tree"`
	Data       map[string][]*t.Transaction `json:"transaction_data"`
	DataMutex  sync.RWMutex
}

// ? Merkle tree simple representation
type Merkle struct {
	Root  *Merkle `json:"root"`
	Hash  string  `json:"hash"`
	Left  *Merkle `json:"left"`
	Right *Merkle `json:"right"`
}

// ? Hash pair struct for later hashing
type HashPair struct {
	LeftHash  string `json:"left_hash"`
	RightHash string `json:"right_hash"`
}

func (m *Merkle) String() string {
	value := fmt.Sprintf("{left: %p,right:%p,root:%p,hash:%v}", m.Left, m.Right, m.Root, []byte(m.Hash))
	return value
}

// ? Create a hash with a hash pair
func (hashPair HashPair) HashThePair() (string, error) {
	encoded, err := json.Marshal(hashPair)
	if err != nil {
		return "", nil
	}
	hashed := sha256.Sum256(encoded)
	return string(hashed[:]), nil
}

// ? Function that build a merkle tree and returns a string and a error
func (db *Database) BuildMerkle(transactions []*t.Transaction) (string, error) {
	if len(transactions) == 0 {
		return "", fmt.Errorf("transactions cannot be empty")
	}
	//? the init nodes
	nodes := []*Merkle{}
	//? Intantiate the first uppon work merkle
	for i := range transactions {
		newMerkle := &Merkle{}
		hash, err := transactions[i].HashTransaction()
		if err != nil {
			return "", err
		}
		newMerkle.Hash = hash
		nodes = append(nodes, newMerkle)
	}
	//? build up the tree
	root := &Merkle{}
	var err error
	if len(nodes) > 1 {
		root, err = buildIntermediate(nodes)
		if err != nil {
			return "", err
		}
	} else {
		//? Hash only that hash
		hashPair := HashPair{
			LeftHash:  nodes[0].Hash,
			RightHash: nodes[0].Hash,
		}
		hash, err := hashPair.HashThePair()
		if err != nil {
			return "", nil
		}
		root = &Merkle{
			Left:  nodes[0],
			Right: nodes[0],
			Hash:  hash,
		}
	}
	//? Add the transactions and also the merkle tree to the database
	db.DataMutex.Lock()
	//? add the transactions to the db
	db.Data[root.Hash] = transactions
	//? add the root merkle
	db.MerkleTree[root.Hash] = root
	db.DataMutex.Unlock()
	//? In the final we can perfectly return the recentMerkles[1].Hash, which represents the root after all the work done on the loop
	return root.Hash, nil
}

// ? Function to build the tree up
func buildIntermediate(data []*Merkle) (*Merkle, error) {
	if len(data) == 1 {
		return data[0], nil
	}
	var nodes []*Merkle
	for i := 0; i < len(data); i += 2 {
		leftIndex, rightIndex := i, i+1
		if i+1 == len(data) {
			rightIndex = i
		}
		//? Create the hash
		pairHash := &HashPair{
			LeftHash:  data[leftIndex].Hash,
			RightHash: data[rightIndex].Hash,
		}
		hash, err := pairHash.HashThePair()
		if err != nil {
			return nil, err
		}
		//? Create a head and put the left and right inside of it
		head := &Merkle{
			Hash:  hash,
			Left:  data[leftIndex],
			Right: data[rightIndex],
		}
		//? Point the left and right to the head
		head.Left.Root = head
		head.Right.Root = head
		nodes = append(nodes, head)
	}
	return buildIntermediate(nodes)
}

// ? Get the pointer with a given hash
func (m *Merkle) getPointerFromHash(hash string) *Merkle {
	//? Stop point
	if m == nil {
		return nil
	}
	if m.Hash == hash {
		return m
	}
	//? Check if it is different than nil
	if leftPointer := m.Left.getPointerFromHash(hash); leftPointer != nil {
		return leftPointer
	}

	if rightPointer := m.Right.getPointerFromHash(hash); rightPointer != nil {
		return rightPointer
	}
	//? return nil
	return nil
}

func (leaf *Merkle) hashUntilReachRoot() (string, error) {
	currentPointer := leaf
	resultHash := ""
	for {
		//? break point
		if currentPointer.Root == nil {
			break
		}
		//? establish the current pointer
		currentPointer = currentPointer.Root
		//? make a hash from it
		hashP := &HashPair{
			LeftHash:  currentPointer.Left.Hash,
			RightHash: currentPointer.Right.Hash,
		}
		hash, err := hashP.HashThePair()
		if err != nil {
			return "", err
		}
		resultHash = hash
		//? Check if the hash is equal to the hash of the current pointer whish is the root of the leaf
		if hash != currentPointer.Hash {
			return "", fmt.Errorf("the result hash is not equal to the hash of the root of the leaf")
		}
	}
	return resultHash, nil
}

// ? Verify hash in a given merkle tree
func (m *Merkle) VerifyHash(hash string) (bool, error) {
	//? get the pointer
	pointer := m.getPointerFromHash(hash)
	if pointer == nil {
		return false, nil
	}
	//? Hash the tree from the leaf until reach the root
	newHash, _ := pointer.hashUntilReachRoot()
	//? check equality
	return newHash != "", nil
}
