package test

import (
	"context"
	"fmt"
	server "raft/protofiles"
	"raft/server/utils"
	"testing"
)

func Test_Find_Delete_If_Needed(t *testing.T) {
	//? Construct a valid append
	leaderCommit := 1
	newEntries := &server.Entries{}
	newEntries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 1
	entrie.Term = 1
	entrie.Log = &server.Log{}
	entrie.Log.Command = "test"
	entrie.Log.Args = []string{"1"}
	newEntries.Entrie = append(newEntries.Entrie, entrie)
	err := ServerImpl.FindAndDeleteIfNeeded(newEntries, int32(leaderCommit))
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if len(State.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("the length should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].IndexOfLog != 1 {
		t.Fatalf("index of log should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].Term != 1 {
		t.Fatalf("the term should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].Log.Command != "test" {
		t.Fatalf("the command of the log should be test")
	}
	if len(State.PersistentState.Entries.Entrie[0].Log.Args) != 1 {
		t.Fatalf("should contain one argument inside")
	}
	if State.PersistentState.Entries.Entrie[0].Log.Args[0] != "1" {
		t.Fatalf("the argument should be '1'")
	}
	//? Establish a commit index that is more high than the one that we have
	leaderCommit = 3
	err = ServerImpl.FindAndDeleteIfNeeded(newEntries, int32(leaderCommit))
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if State.VolatileState.CommitIndex != 1 {
		t.Fatalf("it should a commit index of 1")
	}
	State.PersistentState.Entries.Entrie = []*server.Entrie{}
	State.VolatileState.CommitIndex = 0
}

func Test_Broadcast(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	defer fakeNetwork.ReleaseResources()
	//? Make a broadcast happen (the first should be ok)
	req := &server.AppendRequest{}
	req.IdLeader = "candidate1"
	req.LeaderCommit = 0
	req.PrevLogIndex = 0
	req.PrevLogTerm = 0
	req.Term = 0
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 1
	entrie.Term = 0
	entrie.Log = &server.Log{}
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	entrie.Log.Command = "test"
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	res, err := fakeNetwork.server1.Broadcast(
		fakeNetwork.client2,
		context.Background(),
		req,
		"candidate1",
		false,
	)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if !res {
		t.Fatalf("should not throw a error")
	}
	if len(fakeNetwork.state2.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("it should contain 1 log in it")
	}
	if fakeNetwork.state2.PersistentState.Entries.Entrie[0].IndexOfLog != 1 {
		t.Fatalf("should have a log which index is 1")
	}
	if fakeNetwork.state2.PersistentState.Entries.Entrie[0].Term != 0 {
		t.Fatalf("log term should be 0")
	}
	if len(fakeNetwork.state2.PersistentState.Entries.Entrie[0].Log.Args) != 1 {
		t.Fatalf("should have 1 arg")
	}
	if fakeNetwork.state2.PersistentState.Entries.Entrie[0].Log.Command != "test" {
		t.Fatalf("should have the right command in it")
	}
	//? Make 2 valid broadcasts happen
	for i := 0; i < 2; i++ {
		req := &server.AppendRequest{}
		req.IdLeader = "candidate1"
		req.LeaderCommit = int32(i + 1)
		req.PrevLogIndex = int32(i + 1)
		req.PrevLogTerm = 0
		req.Term = 0
		req.Entries = &server.Entries{}
		req.Entries.Entrie = []*server.Entrie{}
		entrie := &server.Entrie{}
		entrie.IndexOfLog = int32(i + 2)
		entrie.Term = int32(i + 2)
		entrie.Log = &server.Log{}
		entrie.Log.Args = []string{}
		entrie.Log.Args = append(entrie.Log.Args, "1")
		entrie.Log.Command = "test"
		req.Entries.Entrie = append(req.Entries.Entrie, entrie)
		res, err := fakeNetwork.server1.Broadcast(
			fakeNetwork.client2,
			context.Background(),
			req,
			"candidate1",
			false,
		)
		if err != nil {
			t.Fatalf("should not throw a error at this stage %v", err)
		}
		if !res {
			t.Fatalf("should not throw a error")
		}
	}
	//? Lets create an broadcast that should overwrite the log
	req = &server.AppendRequest{}
	req.IdLeader = "candidate1"
	req.LeaderCommit = 2
	req.PrevLogIndex = 2
	req.PrevLogTerm = 0
	req.Term = 0
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie = &server.Entrie{}
	entrie.IndexOfLog = 2
	entrie.Term = 1
	entrie.Log = &server.Log{}
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	entrie.Log.Command = "testas"
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	res, err = fakeNetwork.server1.Broadcast(
		fakeNetwork.client2,
		context.Background(),
		req,
		"candidate1",
		false,
	)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if !res {
		t.Fatalf("should not throw a error")
	}
	if fakeNetwork.state2.PersistentState.Entries.Entrie[1].Log.Command != "testas" {
		t.Fatalf("should have overwrite it")
	}
	//? Lets create a request that tries to send all the content when someone does not have anything (send the data to server3)
	fakeNetwork.state1.Entries = fakeNetwork.state2.Entries
	req = &server.AppendRequest{}
	req.IdLeader = "candidate1"
	req.LeaderCommit = 2
	req.PrevLogIndex = 2
	req.PrevLogTerm = 2
	req.Term = 0
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie = &server.Entrie{}
	entrie.IndexOfLog = 3
	entrie.Term = 0
	entrie.Log = &server.Log{}
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	entrie.Log.Command = "testas"
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	res, err = fakeNetwork.server1.Broadcast(
		fakeNetwork.client3,
		context.Background(),
		req,
		"candidate1",
		false,
	)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if !res {
		t.Fatalf("should not throw a error")
	}
	if fakeNetwork.state3.PersistentState.Entries.Entrie[1].Log.Command != "testas" {
		t.Fatalf("should have overwrite it")
	}
}

func Test_BroadcastAll(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	defer fakeNetwork.ReleaseResources()
	//? Create a log state for the leader
	newEntrie := []*server.Entrie{}
	for i := 0; i < 3; i++ {
		entrie := &server.Entrie{}
		entrie.IndexOfLog = int32(i)
		entrie.Term = 0
		entrie.Log = &server.Log{}
		entrie.Log.Args = []string{}
		entrie.Log.Args = append(entrie.Log.Args, "1")
		entrie.Log.Command = fmt.Sprintf("teste%v", i+1)
		newEntrie = append(newEntrie, entrie)
	}
	fakeNetwork.state1.Entries = &server.Entries{}
	fakeNetwork.state1.Entries.Entrie = newEntrie
	req := &server.AppendRequest{}
	req = &server.AppendRequest{}
	req.IdLeader = "candidate1"
	req.LeaderCommit = 0
	req.PrevLogIndex = 3
	req.PrevLogTerm = 0
	req.Term = 0
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 4
	entrie.Term = 0
	entrie.Log = &server.Log{}
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	entrie.Log.Command = "testas"
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	number, err := fakeNetwork.server1.BroadCastAll(req)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if number != 3 {
		t.Fatalf("should have 3 sends")
	}
	if len(fakeNetwork.state2.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
	if len(fakeNetwork.state3.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
	if len(fakeNetwork.state4.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
}

func Test_Redirect_Leader(t *testing.T) {
	fakeNetwork := InitFake4Servers()
	fakeNetwork.state1.PersistentState.ServerMemberState = utils.Leader
	fakeNetwork.state2.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state3.PersistentState.ServerMemberState = utils.Follower
	fakeNetwork.state4.PersistentState.ServerMemberState = utils.Follower
	defer fakeNetwork.ReleaseResources()
	req := &server.AppendRequest{}
	req = &server.AppendRequest{}
	req.IdLeader = "candidate1"
	req.LeaderCommit = 0
	req.PrevLogIndex = 3
	req.PrevLogTerm = 0
	req.Term = 0
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 4
	entrie.Term = 0
	entrie.Log = &server.Log{}
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	entrie.Log.Command = "testas"
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	//? update the leader id of the server2
	fakeNetwork.state2.PersistentState.LeaderId = "candidate1"
	res, err := fakeNetwork.server2.RedirectToLeader(req)
	if err != nil {
		t.Fatalf("should not throw a error at this stage")
	}
	if !res.Success {
		t.Fatalf("should throw sucess")
	}
	if len(fakeNetwork.state1.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
	if len(fakeNetwork.state2.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
	if len(fakeNetwork.state3.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
	if len(fakeNetwork.state4.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("should contain 1 log")
	}
}

func Test_Log_Append(t *testing.T) {
	req := &server.AppendRequest{}
	req.Term = 1
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 1
	entrie.Term = 1
	entrie.Log = &server.Log{}
	entrie.Log.Command = "test"
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	req.LeaderCommit = 1
	req.IdLeader = "candidate2"
	req.PrevLogIndex = 0
	req.PrevLogTerm = 0
	confirmation, err := ServerImpl.AppendLogsRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("should not throw a error at this state %v", err)
	}
	if !confirmation.Success {
		t.Fatalf("should throw true")
	}
	if confirmation.Term != 1 {
		t.Fatalf("the term should be zero at this stage")
	}
	if len(State.PersistentState.Entries.Entrie) != 1 {
		t.Fatalf("the length should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].IndexOfLog != 1 {
		t.Fatalf("index of log should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].Term != 1 {
		t.Fatalf("the term should be 1")
	}
	if State.PersistentState.Entries.Entrie[0].Log.Command != "test" {
		t.Fatalf("the command of the log should be test")
	}
	if len(State.PersistentState.Entries.Entrie[0].Log.Args) != 1 {
		t.Fatalf("should contain one argument inside")
	}
	if State.PersistentState.Entries.Entrie[0].Log.Args[0] != "1" {
		t.Fatalf("the argument should be '1'")
	}
}

func Test_Invalid_Log(t *testing.T) {
	req := &server.AppendRequest{}
	req.Term = 2
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie := &server.Entrie{}
	entrie.IndexOfLog = 1
	entrie.Term = 1
	entrie.Log = &server.Log{}
	entrie.Log.Command = "test"
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	req.LeaderCommit = 1
	req.IdLeader = "candidate2"
	req.PrevLogIndex = 0
	req.PrevLogTerm = 0
	ServerImpl.AppendLogsRPC(context.Background(), req)
	req = &server.AppendRequest{}
	req.Term = 1
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie = &server.Entrie{}
	entrie.IndexOfLog = 1
	entrie.Term = 1
	entrie.Log = &server.Log{}
	entrie.Log.Command = "test"
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	req.LeaderCommit = 1
	req.IdLeader = "candidate2"
	req.PrevLogIndex = 0
	req.PrevLogTerm = 0
	confirmation, err := ServerImpl.AppendLogsRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if confirmation.Success {
		t.Fatalf("should not throw sucess")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should be zero the term")
	}
	req = &server.AppendRequest{}
	req.Term = 2
	req.Entries = &server.Entries{}
	req.Entries.Entrie = []*server.Entrie{}
	entrie = &server.Entrie{}
	entrie.IndexOfLog = 2
	entrie.Term = 2
	entrie.Log = &server.Log{}
	entrie.Log.Command = "test"
	entrie.Log.Args = []string{}
	entrie.Log.Args = append(entrie.Log.Args, "1")
	req.Entries.Entrie = append(req.Entries.Entrie, entrie)
	req.LeaderCommit = 1
	req.IdLeader = "candidate2"
	req.PrevLogIndex = 3
	req.PrevLogTerm = 3
	confirmation, err = ServerImpl.AppendLogsRPC(context.Background(), req)
	if err != nil {
		t.Fatalf("should not throw a error at this stage %v", err)
	}
	if confirmation.Success {
		t.Fatalf("should not throw sucess")
	}
	if confirmation.Term != 2 {
		t.Fatalf("should be zero the term")
	}
}
