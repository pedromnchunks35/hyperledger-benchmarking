package body

import (
	"context"
	server "raft/protofiles"
	logs "raft/server/logs"
	vote "raft/server/vote"
)

// ? Body of the raft simple server
type RaftServer struct {
	logs.Logs
	vote.Vote
	//? Implement the normal interface
	server.UnimplementedRaftSimpleServer
}

func (rs RaftServer) AppendLogsRPC(ctx context.Context, req *server.AppendRequest) (*server.AppendLogsConfirmation, error) {
	return rs.Logs.AppendLogsRPC(ctx, req)
}

func (rs RaftServer) RequestVoteRPC(ctx context.Context, req *server.VoteRequest) (*server.VoteConfirmation, error) {
	return rs.Vote.RequestVoteRPC(ctx, req)
}
