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
	//? Maintain the references, this is for dont loose the pointers when clearing the other arrays (the refill)
	maintainReferences := []*Merkle{}
	//? Create the new generated merkle point
	recentMerkles := []*Merkle{}
	//? Create the merkles under work
	upponWorkMerkles := []*Merkle{}
	//? Intantiate the first uppon work merkle
	for i := range transactions {
		newMerkle := &Merkle{}
		hash, err := transactions[i].HashTransaction()
		if err != nil {
			return "", err
		}
		newMerkle.Hash = hash
		upponWorkMerkles = append(upponWorkMerkles, newMerkle)
	}
	//? Start the iteration until there is no more items
	for {
		//? Break point
		upponWorkMerklesLen := len(upponWorkMerkles)
		recentMerklesLen := len(recentMerkles)
		//? Case we only have 1 recent, means that the root got reached
		if recentMerklesLen == 1 && upponWorkMerklesLen == 0 {
			break
		}
		//? also the recent becomes nil
		//? case we have more recents than 1 and also the on work slice is zero, we need to provide the recent as fuel
		//? Also, since this is the first iteraction, we need to create the new pointers, after that only the root needs to be created
		if upponWorkMerklesLen == 0 {
			upponWorkMerkles = recentMerkles
			recentMerkles = []*Merkle{}
			break
		}
		//? Get the pair
		var leftIndex, rightIndex int
		if upponWorkMerklesLen > 1 {
			leftIndex = 0
			rightIndex = 1
		} else {
			leftIndex = 0
			rightIndex = 0
		}
		//? Hash the pair
		hashPair := &HashPair{}
		hashPair.LeftHash = upponWorkMerkles[leftIndex].Hash
		hashPair.RightHash = upponWorkMerkles[rightIndex].Hash
		pairHash, err := hashPair.HashThePair()
		if err != nil {
			return "", err
		}
		//? Create fixed reference to the left
		newMerkleLeft := &Merkle{}
		newMerkleLeft.Hash = upponWorkMerkles[leftIndex].Hash
		newMerkleLeft.Left = upponWorkMerkles[leftIndex].Left
		newMerkleLeft.Right = upponWorkMerkles[leftIndex].Right
		//? Create fixed reference for the right
		newMerkleRight := &Merkle{}
		newMerkleRight.Hash = upponWorkMerkles[rightIndex].Hash
		newMerkleRight.Left = upponWorkMerkles[rightIndex].Left
		newMerkleRight.Right = upponWorkMerkles[rightIndex].Right
		//? Append to the index fund
		maintainReferences = append(maintainReferences, newMerkleLeft)
		maintainReferences = append(maintainReferences, newMerkleRight)
		//? Created fixed reference to the root
		newMerkle := &Merkle{}
		newMerkle.Hash = pairHash
		newMerkle.Left = maintainReferences[len(maintainReferences)-2]
		newMerkle.Right = maintainReferences[len(maintainReferences)-1]
		//? add it to the index fund
		maintainReferences = append(maintainReferences, newMerkle)
		//? Put this reference as the father of left and right on the index fund
		maintainReferences[len(maintainReferences)-3].Root = maintainReferences[len(maintainReferences)-1]
		maintainReferences[len(maintainReferences)-2].Root = maintainReferences[len(maintainReferences)-1]
		//? add the newMerkle to the recentMerkles
		recentMerkles = append(recentMerkles, newMerkle)
		//? Case there is only one merkle on the on work slice, we shall make it empty
		if leftIndex == 0 && rightIndex == 0 {
			upponWorkMerkles = []*Merkle{}
		} else {
			//? case there are atleast 2 items, we should get what is in front of that index
			upponWorkMerkles = upponWorkMerkles[2:]
		}
	}
	fmt.Println(len(maintainReferences))
	if len(recentMerkles) != 1 {
		//? For loop to work with the references that we have
		for {
			//? Break point
			upponWorkMerklesLen := len(upponWorkMerkles)
			recentMerklesLen := len(recentMerkles)
			//? case we have more recents than 1 and also the on work slice is zero, we need to provide the recent as fuel
			//? also the recent becomes nil
			if recentMerklesLen != 1 && upponWorkMerklesLen == 0 {
				upponWorkMerkles = recentMerkles
				recentMerkles = []*Merkle{}
			}
			//? Case we only have 1 recent, means that the root got reached
			if recentMerklesLen == 1 && upponWorkMerklesLen == 0 {
				break
			}
		}
	}
	//? Add the transactions and also the merkle tree to the database
	db.DataMutex.Lock()
	//? add the transactions to the db
	db.Data[recentMerkles[0].Hash] = transactions
	//? add the root merkle
	db.MerkleTree[recentMerkles[0].Hash] = recentMerkles[0]
	db.DataMutex.Unlock()
	//? In the final we can perfectly return the recentMerkles[1].Hash, which represents the root after all the work done on the loop
	return recentMerkles[0].Hash, nil
}
