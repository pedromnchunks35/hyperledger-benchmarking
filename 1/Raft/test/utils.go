package test

import (
	"context"
	"log"
	"net"
	server "raft/protofiles"
	serverbody "raft/server/body"
	state "raft/server/state"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type FakeNetwork struct {
	client1 server.RaftSimpleClient
	state1  *state.State
	server1 *serverbody.RaftServer
	client2 server.RaftSimpleClient
	state2  *state.State
	server2 *serverbody.RaftServer
	client3 server.RaftSimpleClient
	state3  *state.State
	server3 *serverbody.RaftServer
	client4 server.RaftSimpleClient
	state4  *state.State
	server4 *serverbody.RaftServer
	sv1     *grpc.Server
	sv2     *grpc.Server
	sv3     *grpc.Server
	sv4     *grpc.Server
	lis1    *bufconn.Listener
	lis2    *bufconn.Listener
	lis3    *bufconn.Listener
	lis4    *bufconn.Listener
}

func (fakeNetwork *FakeNetwork) ReleaseResources() {
	fakeNetwork.sv1.Stop()
	fakeNetwork.sv2.Stop()
	fakeNetwork.sv3.Stop()
	fakeNetwork.sv4.Stop()
	fakeNetwork.lis1.Close()
	fakeNetwork.lis2.Close()
	fakeNetwork.lis3.Close()
	fakeNetwork.lis4.Close()
}

func InitFake4Servers() (fakeNetwork *FakeNetwork) {
	fakeNetwork = &FakeNetwork{}
	// ? Create one simple connection
	fakeNetwork.lis1 = bufconn.Listen(1920 * 1920)
	fakeNetwork.lis2 = bufconn.Listen(1920 * 1920)
	fakeNetwork.lis3 = bufconn.Listen(1920 * 1920)
	fakeNetwork.lis4 = bufconn.Listen(1920 * 1920)
	// ? Create a fake state for server 1
	fakeNetwork.sv1 = grpc.NewServer()
	fakeNetwork.server1 = &serverbody.RaftServer{}
	fakeNetwork.state1 = state.InitState("candidate1")
	fakeNetwork.server1.SetState(fakeNetwork.state1)
	// ? server2
	fakeNetwork.sv2 = grpc.NewServer()
	fakeNetwork.server2 = &serverbody.RaftServer{}
	fakeNetwork.state2 = state.InitState("candidate2")
	fakeNetwork.server2.SetState(fakeNetwork.state2)
	// ? server3
	fakeNetwork.sv3 = grpc.NewServer()
	fakeNetwork.server3 = &serverbody.RaftServer{}
	fakeNetwork.state3 = state.InitState("candidate3")
	fakeNetwork.server3.SetState(fakeNetwork.state3)
	//? server4
	fakeNetwork.sv4 = grpc.NewServer()
	fakeNetwork.server4 = &serverbody.RaftServer{}
	fakeNetwork.state4 = state.InitState("candidate4")
	fakeNetwork.server4.SetState(fakeNetwork.state4)
	// ? Register the servers
	server.RegisterRaftSimpleServer(fakeNetwork.sv1, fakeNetwork.server1)
	server.RegisterRaftSimpleServer(fakeNetwork.sv2, fakeNetwork.server2)
	server.RegisterRaftSimpleServer(fakeNetwork.sv3, fakeNetwork.server3)
	server.RegisterRaftSimpleServer(fakeNetwork.sv4, fakeNetwork.server4)
	//? Inject clients in each other
	// ? serve them
	go func() {
		err := fakeNetwork.sv1.Serve(fakeNetwork.lis1)
		if err != nil {
			log.Fatalf("something is not right with the listening %v", err)
		}
	}()
	go func() {
		err := fakeNetwork.sv2.Serve(fakeNetwork.lis2)
		if err != nil {
			log.Fatalf("something is not right with the listening %v", err)
		}
	}()
	go func() {
		err := fakeNetwork.sv3.Serve(fakeNetwork.lis3)
		if err != nil {
			log.Fatalf("something is not right with the listening %v", err)
		}
	}()
	go func() {
		err := fakeNetwork.sv4.Serve(fakeNetwork.lis4)
		if err != nil {
			log.Fatalf("something is not right with the listening %v", err)
		}
	}()
	// ? Create one client for each
	clientconn1, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(
			func(ctx context.Context, s string) (net.Conn, error) {
				return fakeNetwork.lis1.Dial()
			},
		),
	)
	if err != nil {
		log.Fatalf("something went wrong with the connection %v", err)
	}
	// ? Create the secound client
	clientconn2, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(
			func(ctx context.Context, s string) (net.Conn, error) {
				return fakeNetwork.lis2.Dial()
			},
		),
	)
	if err != nil {
		log.Fatalf("something went wrong with the connection %v", err)
	}
	// ? Create third connection
	clientconn3, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(
			func(ctx context.Context, s string) (net.Conn, error) {
				return fakeNetwork.lis3.Dial()
			},
		),
	)
	if err != nil {
		log.Fatalf("something went wrong with the connection %v", err)
	}
	// ? Create forth connection
	clientconn4, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(
			func(ctx context.Context, s string) (net.Conn, error) {
				return fakeNetwork.lis4.Dial()
			},
		),
	)
	if err != nil {
		log.Fatalf("something went wrong with the connection %v", err)
	}
	//? Register the clients
	fakeNetwork.client1 = server.NewRaftSimpleClient(clientconn1)
	fakeNetwork.client2 = server.NewRaftSimpleClient(clientconn2)
	fakeNetwork.client3 = server.NewRaftSimpleClient(clientconn3)
	fakeNetwork.client4 = server.NewRaftSimpleClient(clientconn4)
	//? Inject the client in which server
	fakeNetwork.server1.InjectClient("candidate2", fakeNetwork.client2)
	fakeNetwork.server1.InjectClient("candidate3", fakeNetwork.client3)
	fakeNetwork.server1.InjectClient("candidate4", fakeNetwork.client4)

	fakeNetwork.server2.InjectClient("candidate1", fakeNetwork.client1)
	fakeNetwork.server2.InjectClient("candidate3", fakeNetwork.client3)
	fakeNetwork.server2.InjectClient("candidate4", fakeNetwork.client4)

	fakeNetwork.server3.InjectClient("candidate1", fakeNetwork.client1)
	fakeNetwork.server3.InjectClient("candidate2", fakeNetwork.client2)
	fakeNetwork.server3.InjectClient("candidate4", fakeNetwork.client4)

	fakeNetwork.server4.InjectClient("candidate1", fakeNetwork.client1)
	fakeNetwork.server4.InjectClient("candidate2", fakeNetwork.client2)
	fakeNetwork.server4.InjectClient("candidate3", fakeNetwork.client3)
	return
}
