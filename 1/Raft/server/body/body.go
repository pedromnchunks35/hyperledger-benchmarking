package body

import (
	"context"
	server "raft/protofiles"
	hearthbeat "raft/server/hearthbeat"
	logs "raft/server/logs"
	"raft/server/state"
	"raft/server/utils"
	vote "raft/server/vote"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ? Body of the raft simple server
type RaftServer struct {
	logs.Logs
	vote.Vote
	hearthbeat.Hearthbeat
	//? Implement the normal interface
	server.UnimplementedRaftSimpleServer
	//? the state
	state *state.State
}

func (rs *RaftServer) SetState(state *state.State) {
	rs.state = state
	rs.Vote.SetState(state)
	rs.Logs.SetState(state)
	rs.Hearthbeat.SetState(state)
}

func (rs RaftServer) SendHearthBeats() error {
	err := rs.Hearthbeat.SendHearthBeats()
	if err != nil && strings.Contains(err.Error(), "has more term than candidate") {
		rs.Vote.BecomeFollower()
	}
	return nil
}

// ? Listener in case we are the leaders
func (rs RaftServer) InitListener() {
	for {
		if rs.state.PersistentState.ServerMemberState == utils.Leader && rs.state.PersistentState.Debug {
			utils.Log("Sending hearthbeats\n")
			rs.SendHearthBeats()
		}
	}
}

// ?  We can become followers if we send a higher term
func (rs *RaftServer) AppendLogsRPC(ctx context.Context, req *server.AppendRequest) (*server.AppendLogsConfirmation, error) {
	return rs.Logs.AppendLogsRPC(ctx, req)
}

// ? Function to handle votes
func (rs *RaftServer) RequestVoteRPC(ctx context.Context, req *server.VoteRequest) (*server.VoteConfirmation, error) {
	return rs.Vote.RequestVoteRPC(ctx, req)
}

// ? Function to handle hearthbeats
func (rs *RaftServer) HearthBeatRPC(ctx context.Context, req *server.HearthBeatRequest) (*server.HearthBeatConfirmation, error) {
	conf, err := rs.Hearthbeat.HearhBeatRPC(ctx, req)
	rs.state.VolatileState.ContractRenewal = true
	if err != nil {
		if strings.Contains(err.Error(), "there is a leader already") {
			rs.Vote.BecomeFollower()
			return conf, nil
		} else if strings.Contains(err.Error(), "has more term than candidate") {
			rs.Vote.BecomeLeader()
			return conf, err
		}
	}
	return conf, err
}

// ? Function to inject a client
func (rs *RaftServer) InjectClient(candidateId string, client server.RaftSimpleClient) {
	rs.state.PersistentState.ServerClients[candidateId] = client
}

// ? Inject clients
func (rs *RaftServer) InjectClients(clients map[string]string) error {
	for key, value := range clients {
		opts := []grpc.DialOption{}
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
		conn, err := grpc.Dial(value, opts...)
		if err != nil {
			return err
		}
		newClient := server.NewRaftSimpleClient(conn)
		rs.InjectClient(key, newClient)
	}
	return nil
}
