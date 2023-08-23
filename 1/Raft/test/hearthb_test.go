package test

import (
	"context"
	server "raft/protofiles"
	"raft/server/utils"
	"strings"
	"testing"
	"time"
)

func Test_SendHearthBeat(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	defer fakeNetwork.ReleaseResources()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	//? Send a hearthbeat
	ch := make(chan error)
	defer close(ch)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	go fakeNetwork.server1.SendHearthBeat(
		ctx,
		"candidate2",
		fakeNetwork.client2,
		ch,
		cancel,
	)
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("it should not throw a error")
		}
	}
	//? Lets make a higher term for the state1 and make state2 send a hearthbeat
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Millisecond)
	ch2 := make(chan error, 2)
	defer close(ch2)
	fakeNetwork.state1.PersistentState.CurrentTerm = 2
	//? Create 2 hearth beats
	go fakeNetwork.server2.SendHearthBeat(
		ctx,
		"candidate1",
		fakeNetwork.client1,
		ch2,
		cancel,
	)
	go fakeNetwork.server2.SendHearthBeat(
		context.Background(),
		"candidate1",
		fakeNetwork.client1,
		ch2,
		cancel,
	)
	for i := 0; i < 2; i++ {
		if err := <-ch2; !(strings.Contains(err.Error(), "candidate candidate1(host), has more term than candidate candidate2") || strings.Contains(err.Error(), "context canceled")) {
			t.Fatalf("should return error because it has higher term")
		}
	}
}

func Test_Send_Hearth_Beats(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	defer fakeNetwork.ReleaseResources()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	err := fakeNetwork.server1.Hearthbeat.SendHearthBeats()
	if err != nil {
		t.Fatalf("should not throw a error at this stage: %v", err)
	}
	//? Making someone of the network having an higher term and try to send a hearth beat
	fakeNetwork.state2.PersistentState.CurrentTerm = 3
	err = fakeNetwork.server1.Hearthbeat.SendHearthBeats()
	if !strings.Contains(err.Error(), "has more term than candidate") {
		t.Fatalf("should return a has more term error")
	}
	//? Making everyone with less term
	fakeNetwork.state1.PersistentState.CurrentTerm = 10
	err = fakeNetwork.server1.Hearthbeat.SendHearthBeats()
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if fakeNetwork.state2.PersistentState.CurrentTerm != 10 {
		t.Fatalf("should update his term to the highest term")
	}
	if fakeNetwork.state3.PersistentState.CurrentTerm != 10 {
		t.Fatalf("should update his term to the highest term")
	}
	if fakeNetwork.state4.PersistentState.CurrentTerm != 10 {
		t.Fatalf("should update his term to the highest term")
	}
}

func Test_Hearth_Beat_RPC(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	defer fakeNetwork.ReleaseResources()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	//? Case it is a candidate it should inform that there is a leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Candidate
	req := &server.HearthBeatRequest{}
	req.IdCandidate = "candidate1"
	req.Term = 1
	conf, err := fakeNetwork.server2.Hearthbeat.HearhBeatRPC(context.Background(), req)
	if !strings.Contains(err.Error(), "there is a leader already") {
		t.Fatalf("should return that there is a leader already")
	}
	if !conf.Ok {
		t.Fatalf("should return success")
	}
	if req.Term != 1 {
		t.Fatalf("should have the term as 1")
	}
	//? case there is no leader setted or the leader is different
	conf, err = fakeNetwork.server3.Hearthbeat.HearhBeatRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if !conf.Ok {
		t.Fatalf("should return success")
	}
	if req.Term != 1 {
		t.Fatalf("should have the term as 1")
	}
	if fakeNetwork.state3.PersistentState.CurrentTerm != 1 {
		t.Fatalf("should update the term")
	}
	if fakeNetwork.state3.PersistentState.LeaderId != "candidate1" {
		t.Fatalf("the leader should be updated")
	}
	if !fakeNetwork.state3.VolatileState.ContractRenewal {
		t.Fatalf("contract should be updated to true")
	}
	//? sending twice the same heathbeat, this will give the leadership again
	req = &server.HearthBeatRequest{}
	req.IdCandidate = "candidate2"
	req.Term = 0
	fakeNetwork.state1.PersistentState.CurrentTerm = 1
	fakeNetwork.server1.Hearthbeat.HearhBeatRPC(context.Background(), req)
	conf, err = fakeNetwork.server1.Hearthbeat.HearhBeatRPC(context.Background(), req)
	if !strings.Contains(err.Error(), "has more term than candidate") {
		t.Fatalf("it should throw the error saying that he has more term than the candidate")
	}
	if conf.Ok {
		t.Fatalf("should get false in return")
	}
}
