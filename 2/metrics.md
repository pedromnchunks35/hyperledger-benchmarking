# Metrics vs Measures
- Metrics are combination of measures
- Measures is what we measure when benchmarking a system

# Measures
- Transaction execution time (time it takes for a single transaction to execute)
- Block Size (the size of a individual block)
- Transaction Size (payload of a request)
- Peer response time (measure the time it takes a peer to respond to a query or request from a peer or another client)
- Orderer Processing Time (time that a orderer takes to order a set of transactions)
- CPU utilization
- Memory utilization (RAM)
- Disk Space usage (Ledger space)
- Transaction validation time (time it takes to validate a transaction, time that a peer takes to endorse a transaction)
- Network latency between peers
- Number of transactions a peer receives
- Resource consumptions during the smart contract execution
- Consensus Round Trip Time, how much time does it take to establish consensus
- Transaction response time

## Metrics
- Transaction Throughput (TPS)
- Latency, the time it takes to a transaction to be submitted and confirmed on the blockchain
- Consensus Algorithm Overhead, the time a transaction takes to go thought the consensus proccess
- Network scalability, check how does the network scales as the number of participants and transactions increase
- Resource utilization, in which point we reach a given limit of resource utilization in terms of (disk,memory,cpu).
- Transaction Validation Time, time that a peer endures to validate a transaction
- Blockchain Size, track all over them time how does the size of the blockchain grows
- Fault tolerance, measure how much time does a peer needs to come back and what can cause the network to stop functioning
- Data Privacy and Confidentiality, Assess how well Fabric's privacy and confidentiality features are implemented. Measure the effectiveness of private data collections and channels.
- Transaction Rate vs Endorsement Policy Complexity. Analyze how the change of the policies can impact the transaction flow
- Consensus algorithm configuration, changing the batch size, timout settings.. etc
- Measure the integration of fabric with external systems