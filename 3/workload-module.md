# Workload module / configuration
- It is the part of caliper that constructs and submits our TXs
- It is here where we implement the business logic, the user behavior
- SUT client
- It is an api
## The api
```
/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;

```
- This is what we need to write.. we need to write the function with that name returning a function that we will create that implements the [WorkloadModuleInterface](https://github.com/hyperledger/caliper/blob/v0.5.0/packages/caliper-core/lib/worker/workload/workloadModuleInterface.js)
- This interface contain 3 async functions: 
  ```
  -> initializeWorkloadModule(), which is a function that i called by the worker processes with the following arguments: 
  workerIndex(number),
  totalWorkers(number),
  roundIndex(number),
  roundArguments(object),
  sutAdapter(connector),
  sutContex(object)
  -> submitTransaction(), which is a function that is runned every time the worker wants to create a transaction. It is supposed to have a very efficient implementation to keep up with high frequency scheduling settings
  -> cleanupWorkloadModule(), which is to release resources
  ```
## Simple base class
- We can use the class that the package gives us, but we would still need to create the submitTransaction function

## Examples
```

const { WorkloadModuleInterface } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleInterface {
    constructor() {
        super();
        this.workerIndex = -1;
        this.totalWorkers = -1;
        this.roundIndex = -1;
        this.roundArguments = undefined;
        this.sutAdapter = undefined;
        this.sutContext = undefined;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.workerIndex = workerIndex;
        this.totalWorkers = totalWorkers;
        this.roundIndex = roundIndex;
        this.roundArguments = roundArguments;
        this.sutAdapter = sutAdapter;
        this.sutContext = sutContext;
    }

    async submitTransaction() {
        let txArgs = {
            // TX arguments for "mycontract"
        };

        return this.sutAdapter.invokeSmartContract('mycontract', 'v1', txArgs, 30);
    }

    async cleanupWorkloadModule() {
        // NOOP
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
```
Example with the base class:
```
const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    async submitTransaction() {
        let txArgs = {
            // TX arguments for "mycontract"
        };

        return this.sutAdapter.invokeSmartContract('mycontract', 'v1', txArgs, 30);
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
```

## info
- To use your workload module you just need to reference it in the section of the workloads on the bench-config file

## Tips
- [Logging control](logging-control.md)
- [Runtime enviroment](runtime-env.md)