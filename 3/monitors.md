# Resource and Transaction Monitors
- Caliper monitor modules are used to collect utilization resource and transaction statistics during the test execution
- It uses two kinds of monitors:
    - Resource Monitors to collect statistics on resource utilization during benchmarking, with monitoring reset between test ronds
    - Transaction monitors to collect worker transaction statistics and provide conditional dispatch actions
## Resource
- This is the resource monitors
- We can config the type of the monitors under monitos.resource, inside of a array
- Permitted monitors are:
    - "process" monitor, is a monitoring of a named process on the host machine and monitors the resources consumed by the running clients. The retrieve statistics are: [memory(max), memory(avg), CPU(max), CPU(avg), Network I/O, Disc I/O]
    ```
    interval- update time in seconds
    processes- array of instructions
    commands- parent process
    args- the file to be executed using the command
    multiOutput- enables handling of the discovery of multiple processes and may be one of: avg (take the average of processes values), sum (sums all process values)
    Example:

    monitors:
     resource:
    - module: process
        options:
        interval: 3
        processes: [{ command: 'node', arguments: 'caliper.js', multiOutput: 'avg' }]
    ```
    - "docker" monitor, is a monitor oriented to docker containers on the host machine or remote machine. By using docker api, we can retrieve docker container statistics such as: [memory(max), memory(avg), CPU(max), CPU(avg), Network I/O, Disc I/O]
    ```
    interval
    cpuUsageNormalization - to normalize cpu usage in a scale to 100 (true/false), default is false
    containers- the addresses of the containers
        monitors:
            resource:
            - module: docker
            options:
            interval: 5
            cpuUsageNormalization: true
            containers:
            - peer0.org1.example.com
            - http://192.168.1.100:2375/orderer.example.com
    ```
    - "prometheus" monitor enables the retrieval of data from prometheus. This will only report based on the explicit user queries that are issued against to Prometheus. If defined the provision of a prometheus server will cause caliper to by default using the prometheus PushGateway
    ```
    url - the prometheus url
    metrics - what the results must include in terms of keywords (include) and the queries description where we have the description of each query (queries) 
    name - metric name
    query - the query itself to be issued to prometheus server
    step - the timing step size to use within the range query
    label - a string to match the returned query, it is used as identifier when populating the record
    statistic - if multiple values are returned, it will aggregate the results by either making: avg(return average of all values),max(return the maximum value from the values),min(return the minimum value from all values),sum(return the sum of all values)
    multiplier - An optional multiplier that may be used to convert exported metrics into a more convenient values (converting bytes to GB)


        monitors:
    resource:
    - module: prometheus
      options:
        url: "http://localhost:9090"
        metrics:
            include: [dev-.*, couch, peer, orderer]
            queries:
                - name: Endorse Time (s)
                  query: rate(endorser_propsal_duration_sum{chaincode="marbles:v0"}[1m])/rate(endorser_propsal_duration_count{chaincode="marbles:v0"}[1m])
                  step: 1
                  label: instance
                  statistic: avg
                - name: Max Memory (MB)
                  query: sum(container_memory_rss{name=~".+"}) by (name)
                  step: 10
                  label: name
                  statistic: max
                  multiplier: 0.000001
    ```
    - PS: we can also add the option "interval" to refresh every x seconds. Also relative to prometheus we can make basic auth using those flags that we writed down in the tables in the runtime configs
## Transaction
- Monitor that act on the completation of transactions
### Logging
- "logging" transaction module is used to aggragte transaction statistics at completion of a test round
```
monitors:
    transaction:
    - module: logging
```
### Prometheus
- "prometheus" transaction module is used to expose current transaction statistics of all workers to a prometheus server 
- this module exposes the following metrics: caliper_tx_submitted (counter),caliper_tx_finished (counter) and caliper_tx_e2e_latency (histogram)
- Usage of prometheus gateway push
```
pushInterval - ms interval of push
pushUrl - Prometheus push gateway
processMetricCollectInterval - time interval for default metrics collection,enabled when present
defaultLabels- object of key/value pairs to 
histogramBuckets- to edit buckets configs

monitors:
    transaction:
    - module: prometheus-push
      options:
        pushInterval: 5000
        pushUrl: "http://localhost:9091"
```
### Grafana visualization
- Grafana is an analytics platform that may be used to query and visualize metrics collected by Prometheus. 
- Metrics available: caliper_tx_submitted,caliper_tx_finished and caliper_tx_e2e_latency
- Default labels: roundLabel: the current test round label, roundIndex: the current test round index, workerIndex: the zero based worker index that is sending the information 

## Resource Charting
- It has support for 2 following charts: horizontal bar and polar area
  ```
  charting:
  bar:
  - metrics: [all | <sting list>]
  polar:
  - metrics: [all | <sting list>]
  ```
- We can either put the single metric or "all" metrics