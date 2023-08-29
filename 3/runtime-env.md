# RunTime Configuration
- The conjunction of settings that we can override when using caliper
- There are 7 locations where we can overwrite settings: 
  1. Memory
  2. Command line arguments
  3. Enviroment variables
  4. Project-level config files
  5. User-level config file
  6. User-level configuration file
  7. Machine-level configuration file
  8. Fallback/default configuration file

## 1. Memory
- We can have access to configuration settings and even change them by using the configUtil package, which as a functionality to get and also to set
- Example:
  ```
  const { ConfigUtil } = require('@hyperledger/caliper-core');

    // Retrieves a setting for your module, if not set, use some default
    const shouldBeFast = ConfigUtil.get             ('mymodule-performance-shoudbefast', /*default:*/   true);

    if (shouldBeFast) { /* ... */ } else { /* ... */ }
  ```
## 2. Command line
- We can also overwrite runtime settings by using the command line like so:
  ```
  caliper launch manager \
    --caliper-workspace yourworkspace/ \
    --caliper-benchconfig yourconfig.yaml \
    --caliper-networkconfig yournetwork.yaml \
    --mymodule-performance-shoudbefast=true
  ```
- Note that the upper example is the same thing as this, because the upper cases are always converted to lower ones and also the "_" is converted to "-":
  ```
  caliper launch manager \
    --caliper-workspace yourworkspace/ \
    --caliper-benchconfig yourconfig.yaml \
    --caliper-networkconfig yournetwork.yaml \
    --MyModule_Performance_ShoudBeFast=true
  ```
## 3. Enviroment variables
- We can also set the configs using enviroment variable as so:
  ```
  export MYMODULE_PERFORMANCE_SHOULDBEFAST=true

    caliper launch manager \
    --caliper-workspace yourworkspace/ \
    --caliper-benchconfig yourconfig.yaml \
    --caliper-networkconfig yournetwork.yaml
  ```

## 4. Configuration files
- We can also pass the desired overwrite of settings by files
- This is a advantage when we want to change multiple settings.. also, we can make version controls of it
- Example in a yaml file
  ```
  mymodule:
  performance:
    shouldbefast: true
  ```
### Project level configs
- The settings that are always overwrited should be in a project level config
- In order to pass such file we should create a yaml file and pass it with the flag --caliper-networkconfig (path of the file)
- Note that you can either pass it in the launch or you can set a enviroment var config config "export CALIPER_PROJECTCONFIG=(path of the file)"
- These settings will overwrite the next settings, this is the master setting

### User-level configs
- If you are overwriting the settings for multiple projects, you should put in a user-level config file
- In order to create user-level configs we either pass --Caliper-UserConfig (path) or we create a enviroment variable "export CALIPER_USERCONFIG=~/.config/my-caliper-config.yaml" or we can even pass it in the project level settings like so:
  ```
   caliper:
    userconfig: (path)
    # additional settings
  ```
### Machine-level
- If we use a common workstation and want to share common settings for that machine accross caliper projects and users
- we use the --Caliper-MachineConfig (path) or we use the export CALIPER_MACHINECONFIG=(path) or we set in the yaml files that have more importance above
  ```
    caliper:
    machineconfig: (path)
  # additional settings
  ```
# Available Settings
## Basic Settings
|Key|Description|
|--|--|
|caliper-benchconfig|Path to the benchmark configuration file that describes the test worker(s),test rounds and monitors|
|caliper-networkconfig|Path to the network configuration file that contains information required to interact with the SUT.|
|caliper-machineconfig|The file path for the machine-level configuration file. Can be relative to the workspace|
|caliper-projectconfig|The file path for the project-level configuration file. Can be relative to the workspace|
|caliper-userconfig|The file path for the user-level configuration file. Can be relative to the workspace|
|caliper-workspace|Workspace directory that contains all configuration information|
|caliper-progress-reporting-enabled|Boolean value for enabling transaction completion progress display by the Caliper manager process|
|caliper-progress-reporting-interval|Numeric value used to specify the caliper progress update frequency, in milliseconds|
## Binding Settings
|Key|Description|
|--|--|
|caliper-bind-args|The additional args to pass to the binding(i.e,npm install) command.|
|caliper-bind-cwd|The CWD to use for the binding (i.e, npm isntall) command|
|caliper-bind-file|The path of a custom binding configuration file that will override the default one|
|caliper-bind-sut|The binding specification of the SUT in the (SUT type):(SDK version) format|
## Reporting Settings
|Key|Description|
|--|--|
|caliper-report-charting-hue|The HUE value to construct the chart color scheme form|
|caliper-report-charting-scheme|The color scheme method to use for producing chart colors|
|caliper-report-charting-transparency|The transparency value[0...1] to use for the charts|
|caliper-report-options|The options object to pass to fs.writeFile|
|caliper-report-path|The absolute or workspace-relative path of the generated report file.|
|caliper-report-precision|Precision(significant digits) for the numbers in the report|
## Logging Settings
|Key|Description|
|--|--|
|caliper-logging-formats-allign|Adds a tab delimiter before the messages to align them in the same space|
|caliper-logging-formats-attributeformat-(attribute)|Specifies the formatting string for the log message attribute (attribute)|
|caliper-logging-formats-json|Indicates that the logs should be serialized in JSON format|
|caliper-logging-formats-label|Adds a specified label to every message. Useful for distributed worker scenario|
|caliper-logging-formats-pad|Pads the log level strings to be the same length|
|caliper-logging-formats-timestamp|Adds a timestamp to the messages with the specified format.|
|caliper-logging-formats-colorize-all|Indicates that all log message attributes must be colorized|
|caliper-logging-formats-colorize-colors-(level)|Sets the color for the log messages with level (level)|
|caliper-logging-targets-(target)-enabled|Sets wether the target transport (target) is enabled or disabled|
|caliper-logging-template|Specifies the message structure through placeholders|
## Worker Management Settings
|Key|Description|
|--|--|
|caliper-worker-communication-method|Indicates the type of the communication between the manager and workers|
|caliper-worker-communication-address|The address of the MQTT broker used for distributed worker management|
|caliper-worker-pollinterval|The interval for polling for new available workers,in milliseconds|
|caliper-worker-remote|Indicates whether the workers operate in distributed mode|
## Benchmark phase Settings
|Key|Description|
|--|--|
|caliper-flow-only-end|Indicates whether to only perform the end command script in the network configuration file|
|caliper-flow-only-init|Indicates whether to only perform the init phase of the benchmark|
|caliper-flow-only-install|Indicates whether to only perform the smart contract install phase of the benchmark|
|caliper-flow-only-start|Indicates whether to only perform the start command script in the network configuration file|
|caliper-flow-only-test|Indicates whether to only perform test phase of the benchmark|
|caliper-flow-skip-end|Indicates whether to skip the end command script in the network configuration file|
|caliper-flow-skip-init|Indicates whether to skip the init phase of the benchmark|
|caliper-flow-skip-install|Indicates whether to skip the smart contract install phase of the benchmark|
|caliper-flow-skip-start|Indicates whether to skip the start command script in the network configuration file.|
|caliper-flow-skip-test|Indicates whether to skip the test phase of the benchmark|
## Authentication Settings
|Key|Description|
|--|--|
|caliper-auth-prometheus-username|Basic authentication username to use authenticate with an existing prometheus server|
|caliper-auth-prometheus-password|Basic authentication password to use authenticate with an existing prometheus server|
|caliper-auth-prometheuspush-username|Basic authentication username to use authenticate with an existing prometheus push gateway|
|caliper-auth-prometheuspush-password|Basic authentication password to use authenticate with an existing prometheus push gateway|
## Fabric Connector Settings
|Key|SUT Version|Description|
|--|--|--|
|caliper-fabric-timeout-invokeorquery|All|The default timeout in seconds to use for invoking or querying transactions.Default is 60 seconds.|
|caliper-fabric-gateway-enabled|1.4|Indicates whether to use the fabric gateway-based SDK API for the 1.4 Fabric SUT. Default is false|
|caliper-fabric-gateway-localhost|1.4 Gateway,2.2|Indicates whether to convert discovered endpoints to localhost. Does not apply if discover is set to false in network config. Default is true|
|caliper-fabric-gateway-querystrategy|1.4 Gateway,2.2|Sets the query strategy to use for 2.2 and 1.4 when gateway is enabled. Default is Round Robin|
|caliper-fabric-gateway-eventstrategy|1.4Gateway,2.2|Sets the event strategy to use for 2.2 and 1.4 when gateway is enabled.Default is any in Invoker Organisation|
|caliper-fabric-latencythreshold|1.4|Determines the reported commit time of a transaction based on the given percentage of event sources|
|caliper-fabric-loadbalancing|1.4|Determines how automatic load balancing is applied|
|caliper-fabric-verify-proposalresponse|1.4|Indicates whether to verify the received proposal responses.|
|caliper-fabric-verify-readwritesets|1.4|Indicates whether to verify that the read-write sets returned by the endorsers match|
### Supported Event Strategies [INFO](https://hyperledger.github.io/fabric-sdk-node/release-1.4/module-fabric-network.html#.DefaultEventHandlerStrategies__anchor)
|Strategy|Corresponds to|
|--|--|
|msp_all|MSPID_SCOPE_ALLFORTX|
|msp_any|MSPID_SCOPE_ANYFORTX|
|network_all|NETWORK_SCOPE_ALLFORTX|
|network_any|NETWORK_SCOPE_ANYFORTX|
### Supported Query Strategies [INFO](https://hyperledger.github.io/fabric-sdk-node/release-1.4/module-fabric-network.html#.DefaultQueryHandlerStrategies__anchor)
|Strategy|Corresponds to|
|--|--|
|msp_single|MSPID_SCOPE_SINGLE|
|msp_round_robin|MSPID_SCOPE_ROUND_ROBIN|