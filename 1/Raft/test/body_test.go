package test

import (
	"raft/server/utils"
	"testing"
)

func Test_Send_Hearth_Beats_Body(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	defer fakeNetwork.ReleaseResources()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	//? case the term is higher
	fakeNetwork.state1.PersistentState.CurrentTerm = 1
	err := fakeNetwork.server2.SendHearthBeats()
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if fakeNetwork.state2.ServerMemberState != utils.Follower {
		t.Fatalf("should have become a follower")
	}
	if fakeNetwork.state1.PersistentState.ServerMemberState != utils.Leader {
		t.Fatalf("should have become a leader")
	}
	err = fakeNetwork.server1.SendHearthBeats()
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if fakeNetwork.state2.ServerMemberState != utils.Follower {
		t.Fatalf("should have become a follower")
	}

}
