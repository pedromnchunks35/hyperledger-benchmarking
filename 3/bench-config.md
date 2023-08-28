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
