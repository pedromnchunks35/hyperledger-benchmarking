package test

import (
	"context"
	"fmt"
	hb "raft/server/hearthbeat"
	"raft/server/utils"
	"strings"
	"testing"
	"time"
)

func Test_SendHearthBeat(t *testing.T) {
	fakeGotError := hb.Err{
		Err: false,
	}
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
		&fakeGotError,
	)
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("it should not throw a error")
		}
	}
	//? Lets make a higher term for the state1 and make state2 send a hearthbeat
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	ch2 := make(chan error)
	defer close(ch2)
	fakeNetwork.state1.PersistentState.CurrentTerm = 2
	//? Create 2 hearth beats
	go fakeNetwork.server2.SendHearthBeat(
		ctx,
		"candidate1",
		fakeNetwork.client1,
		ch2,
		cancel,
		&fakeGotError,
	)
	go fakeNetwork.server2.SendHearthBeat(
		context.Background(),
		"candidate1",
		fakeNetwork.client1,
		ch2,
		cancel,
		&fakeGotError,
	)
	//? Breaking spot
	clientsNumber := 0
	for {
		if clientsNumber == 2 {
			break
		}
		select {
		case err := <-ch2:
			if !strings.Contains(err.Error(), "has more term than candidate") {
				t.Fatalf("should return error because it has higher term")
			}
			clientsNumber++
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
	err := fakeNetwork.server1.SendHearthBeats()
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	//? Making someone of the network having an higher term and try to send a hearth beat
	fakeNetwork.state2.PersistentState.CurrentTerm = 3
	err = fakeNetwork.server1.SendHearthBeats()
	fmt.Println(err)
	if !strings.Contains(err.Error(), "has less term than candidate") {
		t.Fatalf("should throw the less term error")
	}
}
