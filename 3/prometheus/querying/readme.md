# PromQL (PROMETHEUS QUERY LANGUAGE)
- Language that lets the user to aggregate time series data in real time
- The result can either be shown as a graph, viewed as tabular data or consumed by external systems
## Data types
- In PromQL there are 4 types of data
- The type of data influence if the data can be graphed or not. instant vector can be directly graphed
### Instant vector
- a set of time series containing a single sample for each time series, all sharing the same timestamp
### Range vector
- a set of time series containing a range of data points over time for each time series
### Scalar
- Simple numeric floating point value
### String
- a simple string value, currently unused
## Literals
- It uses the same logic of especial chars as golang
```
"this is a string"
'these are unescaped: \n \\ \t'
`these are not unescaped: \n ' " \t`
```
## Float literals
Range of chars
```
```
### example
```
23
-2.43
3.4e-9
0x8f
-Inf
NaN
```
## Time series Selectors
### Instant vector selectors
- It allow a selection of a time series and a single sample value for each at a given timestamp
#### Example
This example selects all time series that have the http_requests_total metric name:
```
http_requests_total
```
We can further make filters like this:
```
http_requests_total{job="prometheus",group="canary"}
```
We can also use the following operators:
- =: Select labels that are exactly equal to the provided string
- !=: Select labels that are not equal to the provided string.
- =~: Select labels that regex-match the provided string.
- !~: Select labels that do not regex-match the provided string.
#### Example 2
```
http_requests_total{environment=~"staging|testing|development",method!="GET"}
```