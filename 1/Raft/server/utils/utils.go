package utils

import (
	server "raft/protofiles"
)

// ? Declaring a ENUM
type ServerMemberState int

const (
	Follower ServerMemberState = iota
	Candidate
	Leader
)

// ? Function to find a given entrie, case it exists it returns it, else it returns nil
func FindLog(entries *server.Entries, logIndex int32, logTerm int32) *server.Entrie {
	for _, value := range entries.Entrie {
		if value.IndexOfLog == logIndex && value.Term == logTerm {
			return value
		}
	}
	return nil
}

// ? Function to check if it represents majority
func RepresentsMajority(numberOfVotes int32, numberOfClients int32) bool {
	if numberOfClients > int32(int(float64(numberOfClients)/2)) {
		return true
	} else {
		return false
	}
}
