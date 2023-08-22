package state

import (
	server "raft/protofiles"
	"raft/server/utils"
)

// ? All of the state
type State struct {
	PersistentState
	VolatileState
	JustLeaderVolatileState
}

// ? The persistent state
type PersistentState struct {
	//? Candidate id
	CandidateId string `json:"id_candidate"`
	//? Current term that we are on
	CurrentTerm int32 `json:"current_term"`
	//? The id of the candidate that i just voted in the election
	MyVote string `json:"my_vote"`
	//? The logs that i am receiving (we will create this object in the protoc)
	Entries *server.Entries `json:"entries"`
	//? The server state in terms of membership (leader,candidate,follower)
	ServerMemberState utils.ServerMemberState `json:"server_member_state"`
	//? "Clients that are servers"
	ServerClients map[string]server.RaftSimpleClient `json:"server_clients"`
	//? Number of gathered votes
	GatheredVotes int32 `json:"gathered_votes"`
	//? Leader id
	LeaderId string `json:"id_leader"`
	//? If it is in debug mode
	Debug bool `json:"debug"`
}

// ? Volatile state
type VolatileState struct {
	//? Highest index of the logs that got commited
	CommitIndex int32 `json:"commit_index"`
	//? Last applied log to the state machine (getting executed)
	LastApplied int32 `json:"last_executed_index"`
	//? The "contract renewal", which is if the leader is renewing the contract with hearthbeats
	ContractRenewal bool `json:"contract_renewal"`
}

// ? Volatile state just for leaders
type JustLeaderVolatileState struct {
	//? Array to specify the next index of log to send of each server
	NextIndexServers map[string]int32 `json:"next_index_servers"`
	//? Array to specify the last index replicated by the leader to the server
	AlreadyReplicatedIndexServers map[string]int32 `json:"already_replicated_index_servers"`
}

func InitState(candidateId string) *State {
	newState := &State{}
	//? Persistent state of raft protocol
	newState.PersistentState.CandidateId = candidateId
	newState.PersistentState.CurrentTerm = int32(0)
	newState.PersistentState.Entries = &server.Entries{}
	newState.PersistentState.Entries.Entrie = []*server.Entrie{}
	newState.PersistentState.MyVote = ""
	newState.PersistentState.ServerMemberState = utils.Follower
	newState.PersistentState.ServerClients = make(map[string]server.RaftSimpleClient)
	newState.PersistentState.Debug = true
	newState.PersistentState.LeaderId = ""
	//? Volatile state from the raft protocol for all servers
	newState.VolatileState.CommitIndex = int32(0)
	newState.VolatileState.LastApplied = int32(0)
	newState.VolatileState.ContractRenewal = true
	//? Volatile state of the leader
	newState.JustLeaderVolatileState.AlreadyReplicatedIndexServers = make(map[string]int32)
	newState.JustLeaderVolatileState.NextIndexServers = make(map[string]int32)
	return newState
}
