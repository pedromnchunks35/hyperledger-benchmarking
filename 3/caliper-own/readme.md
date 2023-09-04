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