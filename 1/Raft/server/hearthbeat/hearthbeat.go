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

type Err struct {
	Err bool
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
func (hb *Hearthbeat) SendHearthBeat(ctx context.Context, key string, value server.RaftSimpleClient, ch chan error, cancel context.CancelFunc, gotError *Err) {
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
	if !gotError.Err {
		ch <- err
	}
}

// ? Send hearthbeats
func (hb *Hearthbeat) SendHearthBeats() error {
	gotError := Err{}
	gotError.Err = false
	//? Send concorrently hearthbeats for everyone
	errorChannel := make(chan error)
	defer close(errorChannel)
	for key, value := range hb.state.PersistentState.ServerClients {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		go hb.SendHearthBeat(ctx, key, value, errorChannel, cancel, &gotError)
	}
	//? Create var for stopping point of channel
	clientsNumber := 0
	for {
		if clientsNumber == len(hb.state.PersistentState.ServerClients) {
			break
		}
		select {
		case err := <-errorChannel:
			if err != nil {
				gotError.Err = true
				return err
			}
			clientsNumber++
		}
	}
	return nil
}

// ? Function to handle hearthbeats
func (hb *Hearthbeat) HearhBeatRPC(ctx context.Context, req *server.HearthBeatRequest) (*server.HearthBeatConfirmation, error) {
	if req.Term > hb.state.PersistentState.CurrentTerm {
		hb.state.PersistentState.CurrentTerm = req.Term
		return answerHearhBeat(false, hb.state.CurrentTerm), fmt.Errorf(
			"candidate %v has less term than candidate %v (host)", hb.state.PersistentState.CandidateId, req.IdCandidate,
		)
	}
	if hb.state.PersistentState.CurrentTerm > req.Term {
		return answerHearhBeat(false, hb.state.CurrentTerm),
			fmt.Errorf("candidate %v(host), has more term than candidate %v",
				hb.state.PersistentState.CandidateId, req.IdCandidate)
	}
	if req.IdCandidate != hb.state.PersistentState.LeaderId {
		hb.state.PersistentState.LeaderId = req.IdCandidate
	}
	return answerHearhBeat(true, hb.state.CurrentTerm), nil
}
