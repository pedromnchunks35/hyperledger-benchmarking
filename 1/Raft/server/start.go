package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	server "raft/protofiles"
	body "raft/server/body"
	state "raft/server/state"

	"google.golang.org/grpc"
)

var (
	PORT      = flag.Int("port", 2000, "the port of the server --port")
	candidate = flag.String("candidate", "", "the id of the candidate --candidate")
	config    = flag.String("config", "", "the path for the config file")
)

func main() {
	flag.Parse()
	bind, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", *PORT))
	if err != nil {
		log.Fatalf("something went wrong listening %v", err)
	}
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	serverImpl := &body.RaftServer{}
	newState := state.InitState(*candidate)
	serverImpl.SetState(newState)
	err = serverImpl.InjectFromJson(*config)
	if err != nil {
		log.Fatalf("error getting the clients from the config json file %v", err)
	}
	server.RegisterRaftSimpleServer(grpcServer, serverImpl)
	log.Printf("server is listening at %v", bind.Addr())
	go serverImpl.Vote.BecomeFollower()
	go serverImpl.InitListener()
	if err := grpcServer.Serve(bind); err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}
}
