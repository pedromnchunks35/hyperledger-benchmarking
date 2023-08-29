# Bench-config constitution
- It is formed by two main parts: Benchmark test settings and Monitoring settings
- The config file can be either JSON or YAML
## Benchmark test settings
|Attibute|Description|
|:--|--:|
|test.name|Short name of the benchmark to show in the report|
|test.description|Detailed description of the benchmark|
|test.workers|Object of worker-related configurations|
|test.workers.number|Specifies the number of worker processes to use for executing the workload|
|test.rounds|Array of objects, each describing the settings of a given round|
|test.rounds[i].label|Short name for a given round, usually it corresponds to the type of Transactions|
|test.rounds[i].txNumber|The number of Transactions Caliper should submit during the round|
|test.rounds[i].txDuration|The length of the round in secouds during which Caliper will submit Transactions|
|test.rounds[i].rateControl|The object describing the [rate controller](rate-controller.md) to use for the round|
|test.rounds[i].workload|The object describing the [workload module](workload-module.md) used for the round|
|test.rounds[i].workload.module|The path to the benchmark workload module that will construct the Transactions to submit|
|test.rounds[i].workload.arguments|Arbitrary object that will be passed to the workload module as config|
### Monitoring Settings
- The monitoring configuration determines what kind of metrics the manager processs can gather and from where. The configuration resides under the monitors attribute.
- [Monitors](./monitors.md)
### Example of full config file
```
test:
  workers:
    number: 5
  rounds:
    - label: init
      txNumber: 500
      rateControl:
        type: fixed-rate
        opts:
          tps: 25
      workload:
        module: benchmarks/samples/fabric/marbles/init.js
    - label: query
      txDuration: 60
      rateControl:
        type: fixed-rate
        opts:
          tps: 5
      workload:
        module: benchmarks/samples/fabric/marbles/query.js
monitors:
  transaction:
  - module: prometheus
  resource:
  - module: docker
    options:
      interval: 1
      containers: ['all']
  - module: prometheus
    options:
      url: "http://prometheus:9090"
      metrics:
        include: [dev-.*, couch, peer, orderer]
        queries:
        - name: Endorse Time (s)
          query: rate(endorser_propsal_duration_sum{chaincode="marbles:v0"}[5m])/rate(endorser_propsal_duration_count{chaincode="marbles:v0"}[5m])
          step: 1
          label: instance
          statistic: avg
```
