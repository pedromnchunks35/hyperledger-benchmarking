# Scenario [NORMAL SCENARIO](The one that we been using for testing)
![Benchmark scenario](assets/BenchMarking-Scenario.drawio.png)
# Possible optimizations in a hyper ledger fabric network
    1. Network Configuration: [ELIGIBLE]
        Consensus Mechanism: Depending on your network's requirements, you can optimize the consensus mechanism (e.g., Raft, Kafka) and its configuration settings to achieve better transaction throughput and latency.

    2. Peers and Ordering Service: [ELIGIBLE]
        Resource Allocation: Ensure that each peer and orderer node has sufficient CPU, memory, and storage resources allocated to handle the expected transaction load efficiently.
        Scaling: Consider adding more peers and orderers as needed to distribute the processing load and enhance network capacity.
        Load Balancing: Implement load balancing to evenly distribute transaction requests across multiple peers or endorsing peers.

    3. Database Configuration: [ELIGIBLE]
        Database Tuning: Optimize the database configuration, indexing, and caching mechanisms to improve query performance and data retrieval.
        Database Scalability: Ensure that your database can scale horizontally as the network grows, allowing for efficient data storage and retrieval.

    4. Caching Mechanisms: [NOT ELIGIBLE(this requires a extra layer, not appliable now)]
        Caching: Implement caching mechanisms to reduce redundant queries to the ledger or to store frequently accessed data, reducing the load on the blockchain.

    5. Channel Configuration: [NOT ELIGIBLE (The benchmarking will be done in a small network for now)]
        Use of Channels: Utilize channels effectively to segregate and isolate data and transactions for different participants or use cases. Proper channel design can improve privacy and scalability.

    6. Endorsement Policies [NOT ELIGIBLE (The benchmarking will be done in a small network for now)]:
        Endorsement Optimization: Carefully design endorsement policies to reduce the number of required endorsements while maintaining security and trust.

    7. Peer Connectivity and Network Topology: [NOT ELIGIBLE (We cannot do such things in the hospital but locally we can but it is not a problem at all for benchmarking)]
        Network Latency: Minimize network latency by optimizing the physical or logical network topology, ensuring peers are well-connected, and using low-latency network links.

    8. Monitoring and Profiling: [NOT ELIGIBLE (This is for production systems)]
        Performance Monitoring: Implement monitoring tools and practices to continuously monitor the network's performance, resource utilization, and transaction latency. Use profiling tools to identify bottlenecks and resource constraints.

    9. Code Efficiency: [NOT ELIGIBLE (We will use a very basic contract)]
        Chaincode Efficiency: Optimize your smart contract (chaincode) code for efficiency and speed. This includes minimizing resource consumption, reducing complex computations, and efficient data storage.
# Important things to notice in the architecture (that we can manipulate of course)
|Concern_Id|Concern|Component|
|---|--|--|
|1|In legacy SDK, the client connects to every peer, collects the endorsements and then communicates directly with the orderer|client/peer|
|2|In the new SDK, the client connects to a single peer and the peer does the rest for him|client/peer|
|3|gateway settings|peer|
|4|keep alive settings|peer|
|5|gossip protocol config|peer|
|6|state configurations|peer|
|7|deliveryclient configurations|peer|
|8|discovery configurations|peer|
|9|limits configurations|peer|
|10|chaincode execution config|peer|
|11|db access configuration|peer|
|12|Request limits|peer|
|13|db indexing|db|
|14|history of updates/operations|db|
|15|Keep alive options|orderer|
|16|Time difference allowed|orderer|
|17|Batch timeout|channel config|
|18|Batch Size|channel config|
|19|Number of workers|client|
|20|Rate of requests|client|
|21|Smart contract function type|client|
# Insights about concerns
|Concern_ID|Term|Definition|Implications/Major Concerns|
|--|--|--|--|
|1|||1. Latency (client waiting for a probable high volume of peers) <br> 2. Network Overhead (client establishing connections with multiple peers and concurrent transactions) <br> 3. Client Resource Utilization|
|2|||1. Latency(peer waiting for a probable high volume of peers)<br>2. Network Overhead (peer waiting for the answer of multiple peers)<br>3. We need to build a connector with this approach for caliper<br>|
|3|endorsementTimeout|The maximum amout of time the client (gateway) will wait for endorsement responses from the endorsing peers|1. Finding balance in the time we intend to apply<br>2. Short can speed up transactions<br>3. To Short can leave incomplete endorsements under heavy load<br>4. To long can congest the network|
|3|broadcastTimeout|The maximum amount of time the client will wait for a transaction to be sucessfully broadcasted to the ordering service|1. Finding balance in the time we intend to apply<br>2. Short can reduce the time for the user to confirm transaction submission<br>3. To short may result in failed transactions if it takes longer because of congestion<br>4. To long may cause congestion|
|3|dialTimeout|The maximum amount of time the client will wait for establishing connections to endorsing peers and the orderer|1. Finding balance in the time we intend to apply<br>2. Short may reduce time to initiate transactions<br>3. Too short may lead to connection failures<br>4. Too long may cause a stuck connection|
|4|interval|Time interval at which a peer sends keepalive messages to other peers|1. Shorter can detect and recover from network issues faster<br>2. To short may result in network overhead|
|4|timeout|How long does a peer waits for a response to a keepalive message before considering the other peer unresponsive|1. Shorter can result in quicker detection of unresponsive peers but may lead to false positives in possible delays<br>2. Longer may be better for delays but it takes more time to identify unresponsive peers|
|4|minInterval|minimum allowable time interval between sucessive keepalive messages sent by a peer to another|1. Short time will produce frequent keepalives, quick detection, reduced latency to detections of unresponsiveness or network disruptions and increase of Overhead<br>2. Long values will result in slow error detection, slower recovery, lower overhead and lower resource usage|
|4|client.timeout && client.interval|control the keepalive interval and timeout of outgoing connections to other peers|1. Health of client connections|
|4|deliveryClient.interval|the time interval at which the delivery client sends keepalive messages to the ordering service to maintain connection|1. Shorter time maintain a responsive and active connection<br>2. Shorter time may increase network traffic|
|4|deliveryClient.timeout|maximum amount of time that the delivery client will wait for a response from the ordering service|1. if delivery client does not receive this response it may think that the connection is now unresponsive or disconnected|
|5|membershipTrackerInterval|||
|5|maxBlockCountToStore|||
|5|maxPropagationBurstLatency|||
|5|maxPropagationBurstSize|||
|5|propagateIterations|||
|5|propagatePeerNum|||
|5|pullInterval|||
|5|pullPeerNum|||
|5|requestStateInfoInterval|||
|5|publishStateInfoInterval|||
|5|stateInfoRetentionInterval|||
|5|publishCertPeriod|||
|5|skipBlockVerification|||
|5|dialTimeout|||
|5|connTimeout|||
|5|recvBuffSize|||
|5|sendBuffSize|||
|5|digestWaitTime|||
|5|requestWaitTime|||
|5|responseWaitTime|||
|5|aliveTimeInterval|||
|5|aliveExpirationTimeout|||
|5|reconnectInterval|||
|5|maxConnectionAttempts|||
|5|msgExpirationFactor|||
|5|election.startupGracePeriod||||
|5|election.membershipSampleInterval|||
|5|election.leaderAliveThreshold|||
|5|election.leaderElectionDuration|||
|5|pvtData.pullRetryThreshold|||
|5|pvtData.transientstoreMaxBlockRetention|||
|5|pvtData.pushAckTimeout|||
|5|pvtData.btlPullMargin|||
|5|pvtData.reconcileBatchSize|||
|5|pvtData.reconcileSleepInterval|||
|6|checkInterval|||
|6|responseTimeout|||
|6|batchSize|||
|6|blockBufferSize|||
|6|maxRetries|||
|7|reconnectTotalTimeThreshold|||
|7|connTimeout|||
|7|reConnectBackoffThreshold|||
|8|authCacheMaxSize|||
|8|authCachePurgeRetentionRation|||
|9|concurrency.endorserService|||
|9|concurrency.deliverService|||
|9|concurrency.gatewayService|||
|9|maxRecvMsgSize|||
|9|maxSendMsgSize|||

# Scenario 1 (1 peer + 1 orderer)
![Scenario 1](assets/Scenario1.drawio.png)
## 1.1.1 
- 2 Workers
- Normal Rate
- Read Function
## 1.1.2
- 4 Workers
- Normal Rate
- Read Function
## 1.1.3
- 8 Workers
- Normal Rate
- Read Function
## 1.2.1
- 2 Workers
- Normal Rate
- Write Function
## 1.2.2
- 4 Workers
- Normal Rate
- Write Function
## 1.2.3
- 8 Workers
- Normal Rate
- Write Function
# Scenario 2 (2 peers + 1 orderer)
![Scenario 2](assets/Scenario2.drawio.png)
## 2.1.1 
- 2 Workers
- Normal Rate
- Read Function
## 2.1.2
- 4 Workers
- Normal Rate
- Read Function
## 2.1.3
- 8 Workers
- Normal Rate
- Read Function
## 2.2.1
- 2 Workers
- Normal Rate
- Write Function
## 2.2.2
- 4 Workers
- Normal Rate
- Write Function
## 2.2.3
- 8 Workers
- Normal Rate
- Write Function
# Scenario 3 (1 peer + 2 orderers)
![Scenario 3](assets/Scenario3.drawio.png)
## 3.1.1 
- 2 Workers
- Normal Rate
- Read Function
## 3.1.2
- 4 Workers
- Normal Rate
- Read Function
## 3.1.3
- 8 Workers
- Normal Rate
- Read Function
## 3.2.1
- 2 Workers
- Normal Rate
- Write Function
## 3.2.2
- 4 Workers
- Normal Rate
- Write Function
## 3.2.3
- 8 Workers
- Normal Rate
- Write Function
# Scenario 4 (2 peers + 2 orderers)
![Scenario 4](assets/Scenario4.drawio.png)
## 4.1.1 
- 2 Workers
- Normal Rate
- Read Function
## 4.1.2
- 4 Workers
- Normal Rate
- Read Function
## 4.1.3
- 8 Workers
- Normal Rate
- Read Function
## 4.2.1
- 2 Workers
- Normal Rate
- Write Function
## 4.2.2
- 4 Workers
- Normal Rate
- Write Function
## 4.2.3
- 8 Workers
- Normal Rate
- Write Function
# Scenario 5 (1 peer org1 + 1 peer org2 + 1 orderer)
![Scenario 5](assets/Scenario5.drawio.png)
## 5.1.1 
- 2 Workers
- Normal Rate
- Read Function
## 5.1.2
- 4 Workers
- Normal Rate
- Read Function
## 5.1.3
- 8 Workers
- Normal Rate
- Read Function
## 5.2.1
- 2 Workers
- Normal Rate
- Write Function
## 5.2.2
- 4 Workers
- Normal Rate
- Write Function
## 5.2.3
- 8 Workers
- Normal Rate
- Write Function
  
[MORE SCENARIOS NEED MORE RESOURCES]