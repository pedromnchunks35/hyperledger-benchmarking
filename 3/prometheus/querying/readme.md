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
#### Bad example
```
{job=~".*"} # Bad!
```
- it is a bad example because it matches all strings with a job,which lead to a broader select that is not desirable
- it is better to besides selecting all jobs also provide another label
```
{job=~".+"}              # Good!
{job=~".*",method="get"} # Good!
```
#### Example
- We can also query by metrics names 
  ```
  {__name__=~"lol:.*"}
  ```
- this will query every metric that has lol: at first of the name metric
- metric name cannot be one of the works:
  - bool
  - on
  - ignoring
  - group_left
  - group_right
 
 ```
    on{} # Bad!
 ```
 ```
  {__name__="on"} # Good!
 ```
### Range Vector Selectors
- Range vector selectos is for time and sits inside of "[]"
  ```
    http_requests_total{job="prometheus"}[5m]
  ```
#### Time durations
- ms
- s
- m
- h
- d
- w
- y
##### Example
```
5h
1h30m
5m
10s
```
### Offset modifier
- It creates a offset to measure in the past relative to the current time
  ```
  http_requests_total offset 5m
  ```
- it needs to follow the selector
  ```
  sum(http_requests_total{method="GET"} offset 5m) // GOOD.
  sum(http_requests_total{method="GET"}) offset 5m // INVALID.
  ```
- We also can measure in the future we can do like this
  ```
  rate(http_requests_total[5m] offset -1w)
  ```
### @ modifier
- it lets us to evaluate the expression at a given timeframe using a unix timestamp
  ```
  http_requests_total @ 1609746000
  ```
- this corresponds to evaluating that expression at 2021-01-04T07:40:00+00:00
- It needs also to follow the selector
  ```
  sum(http_requests_total{method="GET"} @ 1609746000) //GOOD
  sum(http_requests_total{method="GET"}) @ 1609746000 // INVALID.
  ```
- We can use it along with "start()" and "end()"
  ```
  http_requests_total @ start()
  rate(http_requests_total[5m] @ end())
  ```
## Subquery
```
    Syntax: <instant_query> '[' <range> ':' [<resolution>] ']' [ @ <float_literal> ] [ offset <duration> ]
```
- resolution is optional