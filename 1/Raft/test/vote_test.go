package test

import (
	"context"
	server "raft/protofiles"
	"raft/server/utils"
	"testing"
)

func Test_Request_Vote(t *testing.T) {
	req := &server.VoteRequest{}
	req.IdCandidate = "Candidate2"
	req.LastLogIndex = 1
	req.LastLogTerm = 3
	req.Term = 3
	confirmation, err := ServerImpl.RequestVoteRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("not expecting a error at this stage")
	}
	if !confirmation.VoteGranted {
		t.Fatalf("should grant the vote")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should have the term equal to 2")
	}
}

func Test_Request_Vote_Invalid(t *testing.T) {
	//? Get invalid request because it was already voted
	req := &server.VoteRequest{}
	req.IdCandidate = "Candidate2"
	req.LastLogIndex = 1
	req.LastLogTerm = 3
	req.Term = 3
	confirmation, err := ServerImpl.RequestVoteRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("not expecting a error at this stage")
	}
	if confirmation.VoteGranted {
		t.Fatalf("should not grant the vote")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should have the term equal to 2")
	}
	//? Put the vote empty again and test the error of inferior term
	State.PersistentState.MyVote = ""
	//? Get invalid request because it was already voted
	req = &server.VoteRequest{}
	req.IdCandidate = "Candidate2"
	req.LastLogIndex = 1
	req.LastLogTerm = 3
	req.Term = 1
	confirmation, err = ServerImpl.RequestVoteRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("not expecting a error at this stage")
	}
	if confirmation.VoteGranted {
		t.Fatalf("should not grant the vote")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should have the term equal to 2")
	}
	//? Get invalid request because it was the last log incorrect
	req = &server.VoteRequest{}
	req.IdCandidate = "Candidate2"
	req.LastLogIndex = 1
	req.LastLogTerm = 1
	req.Term = 2
	confirmation, err = ServerImpl.RequestVoteRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("not expecting a error at this stage")
	}
	if confirmation.VoteGranted {
		t.Fatalf("should not grant the vote")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should have the term equal to 2")
	}
}

func Test_Random_Time(t *testing.T) {
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.Log("Recovering: %v\n", r)
			}
		}()
		ServerImpl.Vote.RandomTimer()
		ch <- nil
	}()
	if !State.VolatileState.ContractRenewal {
		t.Fatalf("should not be true in this stage, since it must wait between 149 and 170ms")
	}
	<-ch
	if State.VolatileState.ContractRenewal {
		t.Errorf("the contract should have been broken at this stage (it can make the election to fast, maybe not a problem)")
	}
}

func Test_Member_State(t *testing.T) {
	ServerImpl.BecomeLeader()
	if State.PersistentState.ServerMemberState != utils.Leader {
		t.Errorf("should become a leader when invoking this")
	}
	ServerImpl.BecomeFollower()
	if State.PersistentState.ServerMemberState != utils.Follower {
		t.Errorf("should become a follower when invoking this")
	}
	ServerImpl.StartElection()
	if State.PersistentState.ServerMemberState != utils.Candidate {
		t.Errorf("should become a candidate")
	}
	ServerImpl.BecomeFollower()
}
