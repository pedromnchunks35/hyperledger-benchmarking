package vote

import (
	"context"
	"fmt"
	"math/rand"
	server "raft/protofiles"
	state "raft/server/state"
	"raft/server/utils"
	"time"
)

type Vote struct {
	state *state.State
}

// ? Function to set the state
func (vote *Vote) SetState(state *state.State) error {
	if state == nil {
		return fmt.Errorf("invalid pointer")
	}
	vote.state = state
	return nil
}

// ? Create a response for the request vote
func answerVote(answer bool, currentTerm int32) *server.VoteConfirmation {
	return &server.VoteConfirmation{
		Term:        currentTerm,
		VoteGranted: answer,
	}
}

// ? Function that creates a vote request
func (vote *Vote) RequestVoteRPC(ctx context.Context, req *server.VoteRequest) (*server.VoteConfirmation, error) {
	utils.Log("Receiving a vote request from candidate %v\n", req.IdCandidate)
	//? Checking if we already voted for someone
	if vote.state.MyVote != "" {
		return answerVote(false, vote.state.CurrentTerm), nil
	}
	//? Case we did not vote in someone, lets compare the terms
	if vote.state.CurrentTerm > req.Term {
		return answerVote(false, vote.state.CurrentTerm), nil
	}
	//? Compare the state and term of the last log
	if vote.state.VolatileState.CommitIndex > req.LastLogIndex ||
		vote.state.CurrentTerm > req.LastLogTerm {
		return answerVote(false, vote.state.CurrentTerm), nil
	}
	vote.state.MyVote = req.IdCandidate
	return answerVote(true, vote.state.CurrentTerm), nil
}

// ? Func to generate a random timer
func (vote *Vote) RandomTimer() {
	utils.Log("Starting a random timer\n")
	for {
		//? Values that it can rage from
		elegibleValues := make([]int, 30)
		for i := 0; i < len(elegibleValues); i++ {
			elegibleValues[i] = 149 + i
		}
		//? Get one from the array randomly
		randomNum := rand.Intn(len(elegibleValues))
		duration := time.Duration(elegibleValues[randomNum]) * time.Millisecond
		//? Make the thread sleep
		time.Sleep(duration)
		if vote.state.VolatileState.ContractRenewal {
			vote.state.VolatileState.ContractRenewal = false
		} else {
			go vote.StartElection()
			break
		}
	}
}

// ? Function to become a follower
func (vote *Vote) BecomeFollower() {
	utils.Log("Becoming a follower\n")
	vote.state.PersistentState.ServerMemberState = utils.Follower
	go vote.RandomTimer()
}

// ? Function to become a leader
func (vote *Vote) BecomeLeader() {
	utils.Log("Becoming a leader\n")
	vote.state.PersistentState.ServerMemberState = utils.Leader
}

// ? Function to start a election
func (vote *Vote) StartElection() {
	utils.Log("Starting a election\n")
	//? Turn myself into a candidate
	vote.state.PersistentState.ServerMemberState = utils.Candidate
	//? Loop all over the clients that we have and gather their votes
	for key, value := range vote.state.PersistentState.ServerClients {
		//? Create the request
		requestVote := &server.VoteRequest{}
		requestVote.IdCandidate = vote.state.PersistentState.CandidateId
		requestVote.LastLogIndex = vote.state.
			PersistentState.Entries.Entrie[len(vote.state.PersistentState.Entries.Entrie)-1].IndexOfLog
		requestVote.LastLogTerm = vote.state.
			PersistentState.Entries.Entrie[len(vote.state.PersistentState.Entries.Entrie)-1].Term
		requestVote.Term = vote.state.PersistentState.CurrentTerm
		//? Make the call using the client in the map
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		defer cancel()
		confirmation, err := value.RequestVoteRPC(
			ctx,
			requestVote,
		)
		if err != nil {
			utils.ErrorLog("The candidate with key %v does just throw an error: %v \n", key, err)
		}
		//? Case the vote failed and the given term is higher, then we shall become a follower and break this cycle
		if !confirmation.VoteGranted && confirmation.Term > vote.state.CurrentTerm {
			go vote.BecomeFollower()
			vote.state.PersistentState.GatheredVotes = int32(0)
			break
		}
		//? Case the vote got conceided, we shall sum 1 vote
		if confirmation.VoteGranted {
			vote.state.PersistentState.GatheredVotes += int32(1)
		}
	}
	//? Case there is a majority of votes we shall become leader
	if utils.RepresentsMajority(vote.state.GatheredVotes, int32(len(vote.state.PersistentState.ServerClients))) {
		go vote.BecomeLeader()
	}
	//? reset votes
	vote.state.PersistentState.GatheredVotes = int32(0)
}
