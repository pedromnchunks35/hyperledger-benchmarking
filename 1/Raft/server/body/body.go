package body

import (
	s "raft/server/state"
)

// ? Body of the raft simple server
type RaftServer struct {
	//? State representation of raft protocol
	State s.State `json:"state"`
}
