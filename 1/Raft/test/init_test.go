package test

import (
	"raft/server/utils"
	"testing"
)

func Test_Initial_State(t *testing.T) {
	if State.PersistentState.CandidateId != "candidate1" {
		t.Fatalf("the candidateId is wrong and it shouldnt")
	}
	if State.PersistentState.CurrentTerm != 0 {
		t.Fatalf("initial current term should be zero")
	}
	if len(State.PersistentState.Entries.Entrie) != 0 {
		t.Fatalf("initial length of entries should be zero")
	}
	if State.PersistentState.MyVote != "" {
		t.Fatalf("the vote should start empty")
	}
	if State.PersistentState.ServerMemberState != utils.Follower {
		t.Fatalf("it should start as a follower")
	}
	if len(State.PersistentState.ServerClients) != 0 {
		t.Fatalf("it should start with length zero")
	}
	if !State.PersistentState.Debug {
		t.Fatalf("it should be true")
	}
	if State.PersistentState.LeaderId != "" {
		t.Fatalf("it should start as empty string")
	}
	if State.VolatileState.CommitIndex != 0 {
		t.Fatalf("the commited initial index should be zero")
	}
	if State.VolatileState.LastApplied != 0 {
		t.Fatalf("the last applied in machine state index should be zero")
	}
	if !State.VolatileState.ContractRenewal {
		t.Fatalf("the state should start as true")
	}
	if len(State.JustLeaderVolatileState.AlreadyReplicatedIndexServers) != 0 {
		t.Fatalf("the length of the already replicated index servers should be zero")
	}
	if len(State.JustLeaderVolatileState.NextIndexServers) != 0 {
		t.Fatalf("the length of the next index for the server should be zero")
	}
}
