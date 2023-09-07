# Data Model
- Prometheus stores all data in time series
- time series are streams of timestamped values belonging to the same metric and the same set of labeled dimensions
- Prometheus may generate temporary time series from queries
## Metric names and labels
- Time series are identifies by metric name and optional key-value pairs called labels
- Metric name specifies the general feature of a system that is measured (ex: http_requests_total)
- Labels enable prometheus dimensional data model, it identifies a particular dimensional instantiation of that metric (adding or removing a label will create a new time series)
- labels with 2 "_" are reserved for internal use
## Notation
```
<metric name>{<label name>=<label value>,...}
```
- example
  ```
  api_http_requests_total{method="POST",handler="/messages"}

  ```