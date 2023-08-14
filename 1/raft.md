# Raft
[White paper](https://raft.github.io/raft.pdf)
- A consensus algorithm which aim is to be more easier to understand and maintain
- To perform operations the majority of nodes need to have the same configuration
- 3 noticable characteristics from Raft:
 ```
 -> Strong leader: It uses a stronger form of leadership than other consensus algorithms. ex: Logs only flow from the leader to other servers (simplifies log replication management, it is more easy to understand also) 
 -> Leader Election: It is achieved by randomized timers to elect the leaders. By making random elections, it forbiddens other peers of concurrently try to become the leader. Also conflits and events also follow this randomness, making it simplier to solve
 -> Membership changes: In case configuration changing, the majority of nodes dictate the configuration to execute some operations
 ```
- In hyper ledger fabric the term "leader" of the peer takes action when the orderer plays the client role against the peers. The peer leader receives the messages and broadcast them to the other peers
## Replication problem
- It is a concept about replication of state between machine. Multiple machines it identical copies.
- There are a lot of solutions for this problem, one is to have this state in external components such as in kafka (zookeepers). Note that the state is the consensus module which is a bunch of ordered commands that are later executed by a state machine. Inside of a replication solution are multiple states, that were achieved by a certain consensus
### Consensus algorithms properties
- They ensure safety
  ```
  Never return an incorrect result. Under all non byzantine conditions like network delays, partitions, packet losses, duplication and reordering
  ```
- They are fully functional (available)
  ```
  They provide high availability uppon the failure of some of its members
  ```
- They dont depend on timing to ensure the consistency of the logs
  ```
  Faulty clocks can at worst cause availability problems
  ```
- Overall system performance will not be impacted by the minority of nodes
  ```
  Only needs the majority to work perfectly, slower nodes do not take down the performance
  ```