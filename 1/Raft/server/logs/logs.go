package logs

import (
	"context"
	"fmt"
	server "raft/protofiles"
	state "raft/server/state"
	utils "raft/server/utils"
	"time"
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
func (logs *Logs) FindAndDeleteIfNeeded(newEntries *server.Entries, leaderCommit int32) error {
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
		//? Also, overwrite everything in front of it
		if foundInvalid != 0 {
			utils.Log("We just found a invalid log\n")
			logs.state.Entries.Entrie = logs.state.Entries.Entrie[:foundInvalid]
			logs.state.Entries.Entrie = append(logs.state.Entries.Entrie, newEntries.Entrie...)
			break
		} else if foundValid == 0 {
			utils.Log("Joining the log\n")
			//? Case we found a valid one we dont do nothing,
			//? but case is not valid and either invalid,we add it to the entries
			logs.state.Entries.Entrie = append(logs.state.Entries.Entrie, newEntries.Entrie[i])
		}
	}
	//? case the leadercommit is bigger than the commitIndex,
	//? we should update the commitindex to the minimum value between the commit of the leader and the last committed index
	if leaderCommit > logs.state.CommitIndex {
		utils.Log("Updating the commit index\n")
		logs.state.CommitIndex = min(leaderCommit, logs.state.Entries.Entrie[len(logs.state.Entries.Entrie)-1].IndexOfLog)
	}
	return nil
}

// ? Making the broadcast iself
func (logs *Logs) Broadcast(
	client server.RaftSimpleClient,
	ctx context.Context,
	req *server.AppendRequest,
	key string,
	isDone bool,
) (bool, error) {
	if isDone {
		return true, nil
	}
	res, err := client.AppendLogsRPC(ctx, req)
	if err != nil {
		utils.ErrorLog("Something went wrong replicating the log to the server which id is %v", key)
		return false, nil
	}
	//? Case it needs to become a follower
	if res.Term > logs.state.CurrentTerm {
		logs.state.CurrentTerm = res.Term
		return true, fmt.Errorf("someone has a higher term")
	}
	if !res.Success {
		utils.ErrorLog("not sucessfull in sending the request,trying another \n")
		newReq := &server.AppendRequest{}
		newReq.LeaderCommit = req.LeaderCommit
		newReq.IdLeader = req.IdLeader
		newReq.PrevLogIndex = req.PrevLogIndex - int32(1)
		newReq.PrevLogTerm = req.PrevLogIndex - int32(1)
		newReq.Term = req.Term
		newReq.Entries = &server.Entries{}
		index := utils.FindLogByIndex(logs.state.Entries, newReq.PrevLogIndex+int32(1))
		tempSlice := utils.UnionSlices(logs.state.Entries.Entrie[index:], req.Entries.Entrie)
		newReq.Entries.Entrie = tempSlice
		return logs.Broadcast(client, ctx, newReq, key, false)
	} else {
		return logs.Broadcast(client, ctx, req, key, true)
	}
}

// ? Function to broadcast to all the peers
func (logs *Logs) BroadCastAll(req *server.AppendRequest) (int, error) {
	numberOfBroadcasts := 0
	for key, value := range logs.state.ServerClients {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*1,
		)
		defer cancel()
		result, err := logs.Broadcast(value, ctx, req, key, false)
		if err != nil {
			return 0, err
		}
		if result {
			numberOfBroadcasts++
		} else {
			utils.ErrorLog("something went wrong sending a log to the candidate %v", key)
		}
	}
	return numberOfBroadcasts, nil
}

// ? Function to redirect to the leader
func (logs *Logs) RedirectToLeader(req *server.AppendRequest) (*server.AppendLogsConfirmation, error) {
	return logs.state.PersistentState.ServerClients[logs.state.PersistentState.LeaderId].AppendLogsRPC(context.Background(), req)
}

// ? Function to append entries
func (logs *Logs) AppendLogsRPC(ctx context.Context, req *server.AppendRequest) (*server.AppendLogsConfirmation, error) {
	utils.Log("Joining a log\n")
	//? Case the leader was not specified
	if req.IdLeader == "" {
		return logs.RedirectToLeader(req)
	}
	//? Case we are the leaders we need to broadcast to all
	if logs.state.PersistentState.ServerMemberState == utils.Leader {
		//? Add the missing properties
		req.IdLeader = logs.state.PersistentState.CandidateId
		req.LeaderCommit = logs.state.VolatileState.CommitIndex
		if len(logs.state.PersistentState.Entries.Entrie) > 0 {
			req.PrevLogIndex = logs.state.PersistentState.Entries.Entrie[len(logs.state.PersistentState.Entries.Entrie)-1].IndexOfLog
			req.PrevLogTerm = logs.state.PersistentState.Entries.Entrie[len(logs.state.PersistentState.Entries.Entrie)-1].Term
		}
		req.Term = logs.state.CurrentTerm
		res, err := logs.BroadCastAll(req)
		if err != nil {
			return nil, err
		}
		if utils.RepresentsMajority(int32(res), int32(len(logs.state.ServerClients))) {
			//? Commit it in our own
			logs.state.Entries.Entrie = append(logs.state.Entries.Entrie, req.Entries.Entrie...)
		}
	} else {
		//? We should update our term in case the given term is higher
		if logs.state.CurrentTerm < req.Term {
			logs.state.CurrentTerm = req.Term
		}
		//? Case the current term is bigger than the leader term lets return false and return our term
		if logs.state.CurrentTerm > req.Term {
			utils.Log("Our term is higher\n")
			return answerAppend(false, logs.state.CurrentTerm), nil
		}
		//? Case we dont find a log that matches the index and term of the leader prev log, return false
		if (utils.FindLog(
			logs.state.PersistentState.Entries,
			req.PrevLogIndex, req.PrevLogTerm,
		) == nil ||
			(len(logs.state.PersistentState.Entries.Entrie) != 0 &&
				req.PrevLogIndex > logs.state.Entries.Entrie[len(logs.state.PersistentState.Entries.Entrie)-1].IndexOfLog)) &&
			req.PrevLogIndex != int32(0) && req.PrevLogTerm != int32(0) {
			utils.Log("We did not found the lastest log and term given log\n")
			return answerAppend(false, logs.state.CurrentTerm), nil
		}
		err := logs.FindAndDeleteIfNeeded(req.Entries, req.LeaderCommit)
		if err != nil {
			return nil, err
		}
	}
	return answerAppend(true, logs.state.CurrentTerm), nil
}
