package test

import (
	server "raft/protofiles"
	utils "raft/server/utils"
	"testing"
)

func Test_FindLog(t *testing.T) {
	entriesMain := &server.Entries{}
	entries := []*server.Entrie{}
	for i := 0; i < 10; i++ {
		entrie := &server.Entrie{}
		entrie.IndexOfLog = int32(i)
		entrie.Term = int32(i)
		entries = append(entries, entrie)
	}
	entriesMain.Entrie = entries
	//? Getting a valid log
	result := utils.FindLog(entriesMain, int32(1), int32(1))
	if result == nil {
		t.Fatalf("should return a valid result")
	}
	//? Getting a invalid log
	result = utils.FindLog(entriesMain, int32(1), int32(20))
	if result != nil {
		t.Fatalf("should return nil")
	}
}

func Test_FindByIndex(t *testing.T) {
	entriesMain := &server.Entries{}
	entries := []*server.Entrie{}
	for i := 0; i < 10; i++ {
		entrie := &server.Entrie{}
		entrie.IndexOfLog = int32(i)
		entrie.Term = int32(i)
		entries = append(entries, entrie)
	}
	entriesMain.Entrie = entries
	result := utils.FindLogByIndex(entriesMain, int32(2))
	if result != 2 {
		t.Fatalf("it should return the correct index")
	}
	result = utils.FindLogByIndex(entriesMain, int32(40))
	if result != -1 {
		t.Fatalf("should return -1 which stands for not found")
	}
}

func Test_UnionSlices(t *testing.T) {
	slice1 := []*server.Entrie{&server.Entrie{}, &server.Entrie{}}
	slice2 := []*server.Entrie{&server.Entrie{}, &server.Entrie{}, &server.Entrie{}}
	result := utils.UnionSlices(slice1, slice2)
	if len(result) != 5 {
		t.Fatalf("should now contain 5 elements")
	}
}

func Test_Represent_Majority(t *testing.T) {
	result := utils.RepresentsMajority(5, 10)
	if !result {
		t.Fatalf("should throw true")
	}
	result = utils.RepresentsMajority(4, 10)
	if result {
		t.Fatalf("should throw false")
	}
}
