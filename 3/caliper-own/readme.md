# Setting up and Running a Performance Benchmark on an existing Network
## 1. Create Caliper Workspace
- We create a directory
- Install the caliper-cli
```
    npm install --only=prod @hyperledger/caliper-cli@0.5.0
```
- Install the latest fabric SDK
  ```
  npx caliper bind --caliper-bind-sut fabric:2.2
  ```
- Create the following sub-directories: networks,benchmarks and workload
- Create the benchmark files that are responsible for the action that you intend to do, the rounds, the rates and also the workload that you will be doing
- Create a network template file inside of networks. In order to do so you need to do as follows
    - 1. Bring the normal CA crypto material from a given client, which is the client certificate and also the private key. In this example, we putted them inside of a directory called certificates, where the private key will be in the keystore and the certificate in the signcerts
    - 2. Bring the tls ca certificate for tls certificates validation. This one you can put also in certificates and in this exaple we inserted him under a directory called tlscacerts
    - 3. Create a CCP file which stands for (common connection profile), because the normal connector uses the gateway sdk instead of the client sdk. This file will be in the config directory.. you can use it as template.
    - 4. Create the networkconfig file where you will mention how your network is composed
- Create a workload file, which will specify the transaction you will do and the init of it
## 2. Run the benchmark
```
npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/myAssetBenchmark.yaml --caliper-flow-only-test
```
## Notes
- We can send the results from the benchmark via the transaction monitor in the benchmark config
- We can also get resources information from docker by enabling docker daemon
  - To enable docker daemon, we need to edit the file /etc/docker/daemon.json
    ```
    {
    "hosts": ["tcp://0.0.0.0:7531","unix:///var/run/docker.sock"],
    "metrics-addr" : "0.0.0.0:7532"
    }
    ```
  - as you can see we added the hosts ip addresses and then oppened the fireall with firewall-cmd
  - we add the resource monitor in the caliper benchmark config
- another way to make monitoring is to add "metrics-addr" ipaddress and also open the port with firewall-cmd, this becomes a exporter that we config in promtheus target
- another one is to add the cadviser and add it as target also in the prometheus
- later is just the normal steps of using prometheus