package test

import (
	server "raft/protofiles"
	serverbody "raft/server/body"
	state "raft/server/state"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var State *state.State
var ServerImpl *serverbody.RaftServer

func TestMain(m *testing.M) {
	//? Init server implementation
	ServerImpl = &serverbody.RaftServer{}
	//? Init the state of the server
	State = state.InitState("candidate1")
	ServerImpl.Logs.SetState(State)
	ServerImpl.Vote.SetState(State)
	//? Listening
	lis := bufconn.Listen(1920 * 1920)
	defer lis.Close()
	//? Create server
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	//? Register the server
	server.RegisterRaftSimpleServer(grpcServer, ServerImpl)
	m.Run()
}
