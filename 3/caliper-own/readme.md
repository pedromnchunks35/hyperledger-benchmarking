# Setting up and Running a Performance Benchmark on an existing Network
## 1. Create Caliper Workspace
- We create a directory
- Install the caliper-cli
```
    npm install --only=prod @hyperledger/caliper-cli@0.5.0
```
- Create the following sub-directories: networks,benchmarks and workload
- Create a network template file inside of networks