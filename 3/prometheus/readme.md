# Prometheus
- Open source systems monitoring and alerting toolkit 
- Collects and stores its metrics as time series data
## Features
- multi-dimensional data model with time series data identified by metric name and key/value pairs
- PromQL, a flexible query language
- no reliance on distributed storage
- time series collection happens via pull model over HTTP
- pushing time series is supported via an intermediary gateway
- targets are discovered via service discovery or static configuration
- multiple modes of graphing and dashboarding support
## Metrics
- numerical measurements
- time series refers to the recording of changes over time
- what users may want to measure is different from app to app
- It is by this that you know where to make your solution better
## Components
- Prometheus server which stores time series data
- client libraries for app code
- push gateway for short-lived jobs
- special-purpose exporters for services like HAProxy
- an alertmanager for alerts
- various support tools
## Architecture
![Prometheus architecture](../assets/prometheus-architecture.png)
- Prometheus stores metrics from instrumented jobs directly or via intermediary push gateway for short-lived jobs
- Stores data locally and runs rules over this data
# First steps
- Download prometheus (we will not do this because we will get the data via metrics in hyperledger fabrics)
- config prometheus
- starting prometheus
## Concepts
### [Data Model](./datamodel.md)
### [Metric types](./metricTypes.md)
### [Jobs and isntances](./jobs-and-instances.md)
### [Prometheus Remote-Write Specification](./prometheus-remote-write-specification.md)
## Getting started
- Install the rar file of prometheus from the repositories (we should download in the source because if not we will only get the binaries and we want the graphical interface as well)
- Add a job and a target.. for hyper ledger fabric we will add our metrics as so:
  ```
  scrape_configs:
  - job_name: 'fabric-peer'
    static_configs:
      - targets: ['peer1.example.com:9090', 'peer2.example.com:9090']

  - job_name: 'fabric-orderer'
    static_configs:
      - targets: ['orderer1.example.com:9091', 'orderer2.example.com:9091']
  ```