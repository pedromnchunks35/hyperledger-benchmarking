package utils

import (
	server "raft/protofiles"

	"github.com/fatih/color"
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

// ? Function to find a given entrie, case it exists it returns it, else it returns nil
func FindLogByIndex(entries *server.Entries, logIndex int32) int {
	for i, value := range entries.Entrie {
		if value.IndexOfLog == logIndex {
			return i
		}
	}
	return -1
}

// ? Function to sum slices
func UnionSlices(slice1 []*server.Entrie, slice2 []*server.Entrie) []*server.Entrie {
	result := []*server.Entrie{}
	for i := range slice1 {
		result = append(result, slice1[i])
	}
	for i := range slice2 {
		result = append(result, slice2[i])
	}
	return result
}

// ? Function to check if it represents majority
func RepresentsMajority(numberOfVotes int32, numberOfClients int32) bool {
	if numberOfClients+1 > int32(int(float64(numberOfVotes)/2)) {
		return true
	} else {
		return false
	}
}

// ? Function to log a error
func ErrorLog(message string, args ...any) {
	title := color.New(color.FgRed)
	msg := color.New(color.FgHiRed)
	title.Printf("[ERROR] ")
	msg.Printf(message, args...)
}

// ? Function to throw logs
func Log(message string, args ...any) {
	title := color.New(color.FgBlue)
	msg := color.New(color.FgHiBlue)
	title.Printf("[DEBUG] ")
	msg.Printf(message, args...)
}
