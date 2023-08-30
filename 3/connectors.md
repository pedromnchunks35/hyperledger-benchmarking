# Writing connectors
- Most important module in caliper
- It provides abstraction layer between the SUT and the different caliper components
## Requirements for quality connectors
### 1. Keep to the predefined interface
- You must implement the given interface so caliper modules can iteract with the connector
- Do not expose additional capabilities outside of the interface if you are not performance testing a specific SUT
- Make sure your connector behaves similarly to others that followed this guide for users to adapt and experiment with your connector/SUT
### 2. Considere the distributed nature of the SUT
- The connector must be aware of as many SUT nodes as it makes sente to support feature like load balancing or SUT-specific request execution policies
- Hide the network topology as much as you can from other caliper modules
- If you must expose certain nodes to the workload modules, then do that very simply. Dnt expose implementation-specific classes representating the nodes
### 3. Considere de actos in the SUT
- We should concert ourselfs with handling digital identity
- There should be many actos performing different requests
- The connector should be easy to switch between client identities for each request
### 4. Do not reinvent the wheel
- Use the available mature libraries / SDK
### 5. Do not be the bottleneck
- Dont write your own SDK, use what you have
- The connector must be efficient as possible
## Implementing the connector
- In order to implement a connector you need to recognize it as a node.js project
- You have four implementation related tasks:
    1. Implement the connector interface (optionally using the available utility base class)
    2. Implement a factory method for instantiating the connector
    3. Define the schema of your network configuration file
    4. Provide binding configurations for your connector
### 1. Connector interface
- We should add the package "@hyperledger/caliper-core"
- After that we will have access to the interface of the connector
  ```
    class ConnectorInterface extends EventEmitter {
        getType() {}
        getWorkerIndex() {}
        async init(workerInit) {}
        async installSmartContract() {}
        async prepareWorkerArguments(number) {}
        async getContext(roundIndex, args) {}
        async releaseContext() {}
        async sendRequests(requests) {}
    }

    module.exports = ConnectorInterface;

  ```
- We should keep in mind that:
  - The connector is used in two different enviroments: manager and worker processes
  - The connector must expose certain events about the requests, otherwise it's not observable by the Caliper workers, which breaks the scheduling mechanism of Caliper
  - *sendRequests* needs to be implemented carefully and efficienctly
  - The more flexibility you provide, the more features you will have to provide, which makes users happier
#### Interface reference
##### getType
- Description
  ```
  Retrieves a short name for the connector type, usually denoting the SUT, ex: fast-ledger. The name can be used by workload modules capable of targeting multiple types of SUT
  ```
- Return type string
- Returns the name of the connector
#### getWorkerIndex
- Description 
  ```
  Retrieves the zero-based worker process index that instantiated the connector
  ```
- Return type number
- Returns the worker process index
#### init
- Description
  ```
  The method is called by both the manager and (optionally) the worker processes to initialize the connector instance and potentially some aspects of the SUT.
  It is connector specific
  In the manager process it instances one-time init tasks that require iteraction with the SUT such as creating digital identities for example
  In the worker process it processes optional tasks that can be performed later depending of our needs, such as creating data structures
  ```
- Parameters
  ```
  workerInit(boolean), this denotes if it got invoked by a worker or by a manager process
  ```
- Return type Promise
- Returns the promise that will resolve upon method completion
#### installSmartContract
- Description
  ```
  The method is called by the manager process to perform contract deployment on the sut, if allowed remotely
  ```
- Return type Promise
- Returns the promise that will resolve upon method completion
#### prepareWorkerArguments
- Description
  ```
  This method is called by the manager process
  It ensures that the connector instance in the manager process can distribute data to the connector instances in the worker processes
  Perfect place to return for example newly created digital identities to the manager process which in turn will distribute them to the worker process instances for further use
  ```
- Return type Promise(Object[])
- Returns the promise of connector-specific objects for each worker that will resolve upon method completion
#### getContext
- Description
  ```
  The method is called by the worker processes before each round and can be used to assemble a connector-specific object that will be shared with the workload module of the current round.
  The method is also good for get necessary resources for the next round, like establishing connections to remote nodes
  ```
- Parameters
  ```
  roundIndex (number) The zero-based index of the imminent round
  args (object) the object assembled for this worker instance in the prepareWorkerArguments method of the manager instance
  ```
- Return type Promise(object)
- Returns the promise of a connector-specific object that will resolve upon method completion
#### releaseContext
- Description
  ```
  the method is called by the worker processes after each round, and can be used to release resources claimed in the getContext method
  ```
- Return type Promise
- Returns the promise that will resolve upon method completion
#### sendRequests
- Description
  ```
   method called in the worker processes by the workload modules of the rounds.
   Needs to be carefully implemented.
   Must accept one or multiple settings objects that belong to the request that must be sent to the SUT.
   The connector does not have to preserve the order of execution for the requests, unless the SUT type suppoerts such request batches.
   The connector must gather at least the start time, finish time and final status (success or failed) of every request throught TxStatus instances
  ```
- Parameters
  ```
  requests(object or object[]) One or more specific settings object for the request
  ```
- Return type Promise(TxStatus or TxStatus[])
- Return the promise of one or more request execution result that will resolve upon method completion
### Exposed events
- Connector must expose the following events with names matching the defined [constants](https://github.com/hyperledger/caliper/blob/v0.5.0/packages/caliper-core/lib/common/utils/constants.js) for them.
- Without these events the caliper scheduling mechanism wont function well
- Other components also might rely on them (like Tx monitors)
#### txsSubmitted
- Description
  ```
  The event must be reaised when one or more requests are submitted for execution to the SUT. Typically the event should be raised for every individual request
  ```
- Parameters
  ```
  count(number) the number of requests submitted
  ```
#### txsFinished
- Description
  ```
  The event must be raised when one or more requests are fully processed by the SUT (ex: the connector received the results)
  ```
- Parameters retuls (TxStatus or TxStatus[]) One or more request execution result gathered by the connector
### Optional Base Class
- The "@hyperledger/caliper-core" also exports a "ConnectorBase" class that provides sensible default implementations for the connector interface
- This implementations are as below
#### prepareWorkerArguments
- An empty object is returned for each worker by default, which means nothing is shared with the worker process instances
#### sendRequests
- Handles the cases when single or multiple requests are submitted by the workload modules.
- Raises the necessary events before and after the requests.
- The method delegates the execution of a single request to the **_sendSingleRequest** method
#### constructor
- Declares a constructor that requires the worker index and SUT connector type as parameters
#### getType
- Provides a simple getter for the corresponding constructor argument
#### getWorkerIndex
- Provides a simple getter for the corresponding constructor argument
#### More info
- If we opt to use this base class we need to implement the ***_sendSingleRequest*** method
#### _sendSingleRequest
- Description
  ```
  The method only has to handle the sending and processing of a single request
  ```
- Parameters
  ```
  request (object) A connector specific settings object for the request
  ```
- Return type Promise <TxStatus>
- Returns the promise of a request execution result that will resolve upon method completion
### The factory method
- The entry point for your connector implementation will be a factory method
- Manager and worker processes will call this exported factory to instance the connector
#### ConnectorFactory
- Description
  ```
  Instantiates a connector and optionally initializes it. When called from the manager process (denoted with a worker index of -1), the manager will handle calling the init and installSmartContracts methods.
  This initialization is optional in the worker processes, so the factory method must handle it if required
  ```
- Parameters
  ```
  workerIndex(number) The zero-based index of the worker process, or -1 for the manager process
  ```
- Return type Promise(ConnectorInterface)
- Returns the promise of a ConnectorInterface instance that will resolve upon method completion
- Example of possible factory method 
  ```
   'use strict';

    const FastLedgerConnector = require('./fast-ledger-connector');

    async function ConnectorFactory(workerIndex) {
        const connector = new FastLedgerConnector(workerIndex, 'fast-ledger');

        // initialize the connector for the worker processes
        if (workerIndex >= 0) {
            await connector.init(true);
        }

        return connector;
    }

    module.exports.ConnectorFactory = ConnectorFactory;
  ```
### Network configuration file
- Can contain whatever information our connector requires to communicate with the SUT
- It can either be JSON or YAML
- YAML is prefered
- Example of structure
  ```
    # mandatory
    caliper:
        # mandatory
        blockchain: fast-ledger
        # optional
        commands:
            start: startLedger.sh
            end: stopLedger.sh

  ```
- caliper.blockchain attribute tells Caliper which connector to load for the test
- The value depends on how you want to integrate the connector with caliper
### Binding configuration
- The [binding command](https://hyperledger.github.io/caliper/v0.5.0/installing-caliper/#the-bind-command) allows us to install connector dependencies. SUT SDKS and other client libraries usually fall into this category. 
- Example of configuration
  ```
  sut:
  fast-ledger:
    1.0:
      packages: ['fast-ledger-sdk@1.0.0']
    1.4:
      packages: ['fast-ledger-sdk@1.4.5']
    2.0: &fast-ledger-latest
      packages: ['fast-ledger-sdk@2.0.0']
    latest: *fast-ledger-latest
  ```
- Notes about this
  1. Since sut top-level attribute denotes the config section you can and should use anchors to improve the readability in YAML
  2. sut attribute contains keys that identify the SUT types
  3. under the type we can define several SUT versions that our connector supports. It is recommended to use keys corresponding to the semanting version of the SUT. You can specify the type then in the command line using ***--caliper-bind-sut (name):(version)***
  4. Every SUT version needs to declare the required packages caliper should download during the runtime. Different SUT versions must have different SDKS.
  5. Even thought in the example we have version 1.4, we download the version 1.4.5 because it is always good to have the latest version which has the latest bug fixing
  6. If you put the version as latest ex: ***(name):latest***, then when you pass the flag you can only pass the name ***--caliper-bind-sut (name)***, if using a library management like npm and dockerhub ofc
- Another example but more advanced: (case you are using later versions than need extra flags or arguments or different compilers)
  ```
  sut:
  fast-ledger:
    1.0:
      packages: ['fast-ledger-sdk@1.0.0', 'comm-lib@1.0.0']
      settings:
      # compiling older comm-lib on newer Node.js version
      - versionRegexp: '^((?!v8\.).)*$'
        env:
          CXXFLAGS: '-Wno-error=class-memaccess'
          CFLAGS: '-Wno-error=class-memaccess'
        args: '--build-from-source'
  ```
### Documenting the connector
-  Providing proper user manual for your connector is just as important as quality implementation
-  You should make a short summary of it including the supported SUT versions, SUT types and the capabilities of our connector (SUT features and also limitations)
#### Installing dependencies
- Case the connector supports multiple versions, we shall document how do we use such versions
- Document the possible bindings and which bindings have some limitations
#### Runtime settings
- Document how can the end user change some of the behavior of the SUT using the runtime settings
#### Request API
- The main users of the connector will be workload module devs. They will interact with the connector mainly through the sendRequests method. We shall document which settings we accept in our connector.. the settings that will go inside of the sendRequests that will accept either a single or multiple settings objects
- This settings typically include:
  - Operations to execute on the SUT
  - Args of the ops
  - The identity who should submit the request
  - The node(s) to send the request to
  - Differentiation between read-only/write requests
#### Gathered request data
- Connector must report basic execution data towards Caliper to ensure correct reporting
- You are also free to collect any kind of client-side data you may have access to
- We need to know what data users will find useful
- We shall document such collected data (semantics and data types)
#### Network configuration file
- The most important piece of documentation is the schema of the network configuration file your connector can process
- We shall document that structure to define the topology, participants and any required artifacts
- We should also document the semantis and data types of different settings
- Document constraints that could arise between multiple attributes (mutual exclusion,valid values, etc...)
#### Example network configuration
- Provide a fully specified and functioning network config example, because it gots easier for some to absorb a concrete example
### Integration with caliper
- Once a connector is implemented, you can either use it as a 3rd party, pluggable component, which is part of benchmark project
- Contribute your connector to the official caliper code-base, so its always isntalled together with caliper
#### 3rd party connector
- You can easily plug the connector dynamically without it being part of the Caliper code-base. The process is the following:
    1. Create a index.js file in your project that will export your connector factory. The file provides a clean entry point for your connector:
    ```
     'use strict';
        module.exports.ConnectorFactory = require   ('./lib/connectorFactory').ConnectorFactory;
    ```
    2. Set the ***./fast-ledger/index.js*** path for the ***caliper.blockchain*** atribute in your network config file. Path should be relative to the caliper workspace directory or absolute.Caliper will load the module from this path.
    3. If we support different binding, prepare the binding config file for your connector
    4. When you launch caliper, the connector implementation will be picked up through the network config file
    5. You can specify a custom binding config using for example ***--caliper-bind-file ./fast-ledger-binding.yaml*** , not forgetting to specify the binding with ***--caliper-bind-sut fast-ledger:1.0***
    - You can also set the ***caliper.blockchain*** attribute to an NPM package name if you published your connector, in that case the package must be installed in the caliper workspace directory prior to running the benchmark. The convention name must be something like ***caliper-sut*** for the package.
#### Built-in connector
- Contributing with a connector to the code-base, you also accept the responsibility of maintaining the connector when needed. Otherwhise it may become deprecated in future releases
  1. Create a directory with a desired name in the packages directory of the repository for putting the implementation inside
  2. Update your metadata in your own package.json file accordingly. The package name should be scoped ***hyperledger/(name of the directory)***
  3. If your connector supports binding, then you should list dynamic packages in the devDependencies section, so you are not automatically installed with caliper. Also, add the connector binding specifications in the built-in binding config file
  4. Add your new directory path to the root ***lerna.json*** file, under the package section. This will ensure that your package is boostrapped correctly for other developers (and for testing,publishing,etc...)
  5. Add your new package (by name) to the Caliper CLI dependencies
  6. List your connector as a built-in connector in the ***caliper-utils.js*** module,under the BuiltinConnectors variable:
  ```
    const BuiltinConnectors = new Map([
    ['fast-ledger', '@hyperledger/         caliper-fast-ledger'],
    // other connectors...
    ]);
  ```
  7. It is recommended to also provide integration tests for the connector
  8. Make sure the code-related artifact contains the appropriate licence header
  9. You are done, other devs can use your connector