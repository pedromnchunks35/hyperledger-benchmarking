package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	server "raft/protofiles"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	IP string `json:"ip"`
}

var (
	config = flag.String("config", "", ".--config (name of the file config for the connection)")
	logs   = flag.String("logs", "", "--logs (path for the file that contains the logs content)")
)

func main() {
	flag.Parse()
	fileConfig, err := os.ReadFile(*config)
	if err != nil {
		log.Fatalf("%v", err)
	}
	var decodedConfig Config
	err = json.Unmarshal(fileConfig, &decodedConfig)
	if err != nil {
		log.Fatalf("should decode correctly %v", err)
	}
	fileLogs, err := os.ReadFile(*logs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	var decodedLogs []*server.Entrie
	err = json.Unmarshal(fileLogs, &decodedLogs)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(decodedConfig.IP, opts...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()
	client := server.NewRaftSimpleClient(conn)
	req := &server.AppendRequest{}
	req.Entries = &server.Entries{}
	req.Entries.Entrie = decodedLogs
	conf, err := client.AppendLogsRPC(context.Background(), req)
	if err != nil {
		log.Fatalf("should not throw a error %v", err)
	}
	fmt.Println(conf)
}
