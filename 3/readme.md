# Caliper
- Blockchain performance benchmark framework, which allows users to test diffeerent blockchain solutions with custom use cases
## Supported blockchain solutions
- Hyperledger Besu
- Ethereum
- Hyperledger fabric
- FISCO BCOS
## Supported performance metrics
- Transaction/read throughput
- Transaction/read latency
- Resource consumption
## Architecture
### Bird's eye view
- Caliper is a service that generates a workload against a specific system under test (SUT) and continuously monitors its responses. 
- It generates a report based on the observed SUT responses.
![Caliper benchmark](../2/assets/caliper-bench-mark.png)
### Benchmark config file
- This is a file that describes how the benchmark should be executed
- It tells caliper how many rounds it should execute
- At what rate the TXs should be submitted
- Which module will generate the TX content
- It includes settings about monitoring the SUT (which is the system under test)
- Settings are SUT independent, you can use the same settings for multiple versions of your system but this characteristic can be avoided if you specifically target a specific version of the system
- Basicly it is the responsible for dictate the execution of the workloads and the results of the benching
- See the constitution [here](bench-config.md)
### Benchmark artifacts
- Crypto materials necessary to interact with the SUT
- Smart contract source code for caliper to deploy (If the SUT connector support such operation)
- Runtime config files 
- Pre-installed third party packages for your workload modules
### [Writing connectors](connectors.md)