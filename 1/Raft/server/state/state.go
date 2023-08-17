package state

// ? All of the state
type State struct {
	PersistentState
	VolatileState
	JustLeaderVolatileState
}

// ? The persistent state
type PersistentState struct {
	//? Current term that we are on
	CurrentTerm int `json:"current_term"`
	//? The id of the candidate that i just voted in the election
	MyVote string `json:"my_vote"`
	//? The logs that i am receiving (we will create this object in the protoc)
	Log []string `json:"log"`
}

// ? Volatile state
type VolatileState struct {
	//? Highest index of the logs that got commited
	CommitIndex int `json:"commit_index"`
	//? Last applied log to the state machine (getting executed)
	LastApplied int `json:"last_executed_index"`
}

// ? Volatile state just for leaders
type JustLeaderVolatileState struct {
	//? Array to specify the next index of log to send of each server
	NextIndexServers []map[string]int `json:"next_index_servers"`
	//? Array to specify the last index replicated by the leader to the server
	AlreadyReplicatedIndexServers []map[string]int `json:"already_replicated_index_servers"`
}
