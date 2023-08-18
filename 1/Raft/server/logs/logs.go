package logs

import (
	"context"
	"fmt"
	server "raft/protofiles"
	state "raft/server/state"
	utils "raft/server/utils"
)

type Logs struct {
	state *state.State
}

// ? Function to set the state
func (logs *Logs) SetState(state *state.State) error {
	if state == nil {
		return fmt.Errorf("invalid pointer")
	}
	logs.state = state
	return nil
}

// ? Function to asnwer the leader
func answerAppend(success bool, currentTerm int32) *server.AppendLogsConfirmation {
	return &server.AppendLogsConfirmation{
		Term:    currentTerm,
		Success: success,
	}
}

// ? Function to find a given index and in case we find it, delete what is ahead and append the new given item
func (logs Logs) findAndDeleteIfNeeded(newEntries *server.Entries, leaderCommit int32) error {
	//? We loop all over the new entries and try to find them in the server entries
	for i := range newEntries.Entrie {
		foundInvalid := 0
		foundValid := 0
		for j := range logs.state.Entries.Entrie {
			//? Case we find a log with the same index as the leader, we put him as invalid
			if newEntries.Entrie[i].IndexOfLog == logs.state.Entries.Entrie[j].IndexOfLog {
				foundInvalid = j
			}
			//? Case the term of the previous invalid entrie is equal to the leader, we mark him as valid and remove the invalid mark
			if newEntries.Entrie[i].Term == logs.state.Entries.Entrie[j].Term {
				foundInvalid = 0
				foundValid = j
			}
		}
		//? Delete everything in front of that invalid entrie case we have a invalid index
		if foundInvalid != 0 {
			logs.state.Entries.Entrie = logs.state.Entries.Entrie[:foundInvalid]
			logs.state.Entries.Entrie = append(logs.state.Entries.Entrie, newEntries.Entrie...)
		} else if foundValid == 0 {
			//? Case we found a valid one we dont do nothing,
			//? but case is not valid and either invalid,we add it to the entries
			logs.state.Entries.Entrie = append(logs.state.Entries.Entrie, newEntries.Entrie[i])
		}
	}
	//? case the leadercommit is bigger than the commitIndex,
	//? we should update the commitindex to the minimum value between the commit of the leader and the last committed index
	if leaderCommit > logs.state.CommitIndex {
		logs.state.CommitIndex = min(leaderCommit, logs.state.Entries.Entrie[len(logs.state.Entries.Entrie)].IndexOfLog)
	}
	return nil
}

// ? Function to append entries
func (logs *Logs) AppendLogsRPC(ctx context.Context, req *server.AppendRequest) (*server.AppendLogsConfirmation, error) {
	//? Case the current term is bigger than the leader term lets return false and return our term
	if logs.state.CurrentTerm > req.Term {
		return answerAppend(false, logs.state.CurrentTerm), nil
	}
	//? Case we dont find a log that matches the index and term of the leader, return false
	if utils.FindLog(
		logs.state.PersistentState.Entries,
		req.PrevLogIndex, req.PrevLogTerm,
	) == nil {
		return answerAppend(false, logs.state.CurrentTerm), nil
	}
	err := logs.findAndDeleteIfNeeded(req.Entries, req.LeaderCommit)
	if err != nil {
		return nil, err
	}
	return answerAppend(true, logs.state.CurrentTerm), nil
}
