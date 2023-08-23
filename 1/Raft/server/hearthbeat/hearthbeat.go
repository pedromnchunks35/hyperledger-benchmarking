package hearthbeat

import (
	"context"
	"fmt"
	server "raft/protofiles"
	state "raft/server/state"
	"raft/server/utils"
	"time"
)

type Hearthbeat struct {
	state *state.State
}

func (hb *Hearthbeat) SetState(state *state.State) {
	hb.state = state
}

func answerHearhBeat(ok bool, term int32) *server.HearthBeatConfirmation {
	return &server.HearthBeatConfirmation{
		Term: term,
		Ok:   ok,
	}
}

// ? Send hearthbeat
func (hb *Hearthbeat) SendHearthBeat(ctx context.Context, key string, value server.RaftSimpleClient, ch chan error, cancel context.CancelFunc) {
	defer func() {
		if r := recover(); r != nil {
			utils.Log("Recovering from a error in a hearthbeatsend: %v\n", r)
		}
	}()
	defer cancel()
	//? Applying debug mode
	if hb.state.PersistentState.Debug {
		utils.Log("Sending hearthbeat from %v to %v\n", hb.state.PersistentState.CandidateId, key)
	}
	//? Prepare a request and send to everyone
	req := &server.HearthBeatRequest{}
	req.IdCandidate = hb.state.PersistentState.CandidateId
	req.Term = hb.state.PersistentState.CurrentTerm
	_, err := value.HearthBeatRPC(ctx, req)
	ch <- err
}

// ? Send hearthbeats
func (hb *Hearthbeat) SendHearthBeats() error {
	//? Send concorrently hearthbeats for everyone
	errorChannel := make(chan error, len(hb.state.PersistentState.ServerClients))
	defer close(errorChannel)
	for key, value := range hb.state.PersistentState.ServerClients {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		go hb.SendHearthBeat(ctx, key, value, errorChannel, cancel)
	}
	//? Create var for stopping point of channel
	for i := 0; i < len(hb.state.PersistentState.ServerClients); i++ {
		if err := <-errorChannel; err != nil {
			return err
		}
	}
	return nil
}

// ? Function to handle hearthbeats
func (hb *Hearthbeat) HearhBeatRPC(ctx context.Context, req *server.HearthBeatRequest) (*server.HearthBeatConfirmation, error) {
	if hb.state.PersistentState.ServerMemberState == utils.Candidate {
		hb.state.PersistentState.LeaderId = req.IdCandidate
		hb.state.PersistentState.MyVote = ""
		return answerHearhBeat(true, hb.state.CurrentTerm), fmt.Errorf(
			"there is a leader already",
		)
	}
	if req.IdCandidate != hb.state.PersistentState.LeaderId {
		hb.state.PersistentState.LeaderId = req.IdCandidate
		utils.Log("updating the leader\n")
	}
	if req.Term > hb.state.PersistentState.CurrentTerm {
		utils.Log("it has a lower term, lets update it\n")
		hb.state.PersistentState.CurrentTerm = req.Term
	}
	if hb.state.PersistentState.CurrentTerm > req.Term {
		return answerHearhBeat(false, hb.state.CurrentTerm),
			fmt.Errorf("candidate %v(host), has more term than candidate %v",
				hb.state.PersistentState.CandidateId, req.IdCandidate)
	}
	return answerHearhBeat(true, hb.state.CurrentTerm), nil
}
