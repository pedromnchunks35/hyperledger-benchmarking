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