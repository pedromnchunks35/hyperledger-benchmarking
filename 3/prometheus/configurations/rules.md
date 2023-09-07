# Defining recording rules
## Configuring rules
- There are 2 types of rules: recording rules and alerting rules
- To assign this configurations, we have a parameter pointing for that configurations
### Check rules syntax
- In order to check syntax rules without starting a prometheus server all over again we can do:
```
promtool check rules (path to rules)
```
### Recording rules
- This rules allow you to precompute frequently needed or computationally expensive expressions and save their results as a new set of time series
- This is usefull for dashboards that are always repeating the access to a given information
- Recording and alerting rules exist in a rule group
- Rules within a group are run sequentially at a regular interval, with the same evaluation time
- The names of recording rules must be valid metric names
- The names of alerting rules must be valid label values
- The syntax is as follows:
  ```
  groups:
  [ - <rule_group> ]
  ```
  ```
    groups:
  - name: example
    rules:
    - record: code:prometheus_http_requests_total:sum
      expr: sum by (code) (prometheus_http_requests_total)
  ```
#### Rule group
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|name|the name of the group|string||
|interval|How often rules in the group are evaluated|duration|global.evaluation_interval|
|limit|Limit the number of alerts an alerting rule and series a recording rule can produce|int|0,which means no limit|
|rules||array of rule||
#### Rule
##### Record rule
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|record|the name of the time series to output to|string||
|expr|The PromQL expression to be evaluated to, must be a valid metric name|string||
|labels|Labels to add or overwrite before storing the result|labelname:labelvalue||
##### Alerting Rule
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|alert|the name of the alert|string||
|expr|The PromQL expression to evaluate|string||
|for|Time until firing the alert|duration|0s|
|keep_firing_for|How long a alert will continue to fire|duration|0s|
|labels|Overwrite a label with another name for the alert|labelname:newlabel||
|annotations|annotations to add to each alert by labelname|labelname:annotation||
###### Example of template for alert
```
groups:
- name: example
  rules:

  # Alert for any instance that is unreachable for >5 minutes.
  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} down"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  # Alert for any instance that has a median request latency >1s.
  - alert: APIHighRequestLatency
    expr: api_http_request_latencies_second{quantile="0.5"} > 1
    for: 10m
    annotations:
      summary: "High request latency on {{ $labels.instance }}"
      description: "{{ $labels.instance }} has a median request latency above 1s (current value: {{ $value }}s)"
```
## Unit testing Rules
- In order to unit testing the rules we use **promtool**
```
./promtool test rules (TEST FILES, it can be more than 1 separed by spaces)
```
### Test file format
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|rule_files|List of rules files to be tested|array of file_name||
|evaluation_interval|The evaluation interval for the test files|duration|1m|
|group_eval_order|order groups for testing here|array of group_name||
|tests|All the tests are here|array of test_group||
#### test_group
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|interval|interval between series|duration||
|input_series|series config|array of series||
|name|name of the group|string||
|alert_rule_test|unit testing for alert rules|array of alert_test_case||
|promql_expr_test|unit testing for promql expressions|array of prompl_test_case||
|external_labels|external labels acessible to the alert template|labelname:string||
|external_url|external URL acessible to the alert template|string||
#### series
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|series|it follows the usual series notation|string||
|values||string||
#### alert_test_case
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|eval_time|the time for alerts to be checked|duration|0s|
|alertname|the name of the alert to be tested|string||
|exp_alerts|list of expected aleerts which are firing under the given alertname at given evaluation time. case we want to test that the alert does not fire, we put the above fields and leave this empty|array of alert||
#### alert
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|exp_labels|expanded labels of the expected alert|labelname:string||
|exp_annotations|expended annotations from the expected alert|labelname:string||
#### promql_test_case
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|expq|Expression to evaluate|string||
|eval_time|time where expression must be evaluated|duration||
|exp_samples|expected samples at the given evaluation time|array of sample||
##### sample
|Property|Definition|Type/Syntax|Default|
|--|--|--|--|
|labels|labels of the sample in usual series notation|string||
|value|expected value of the promql expression|number||
## Example
```
    # This is the main input for unit testing.
# Only this file is passed as command line argument.

rule_files:
    - alerts.yml

evaluation_interval: 1m

tests:
    # Test 1.
    - interval: 1m
      # Series data.
      input_series:
          - series: 'up{job="prometheus", instance="localhost:9090"}'
            values: '0 0 0 0 0 0 0 0 0 0 0 0 0 0 0'
          - series: 'up{job="node_exporter", instance="localhost:9100"}'
            values: '1+0x6 0 0 0 0 0 0 0 0' # 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0
          - series: 'go_goroutines{job="prometheus", instance="localhost:9090"}'
            values: '10+10x2 30+20x5' # 10 20 30 30 50 70 90 110 130
          - series: 'go_goroutines{job="node_exporter", instance="localhost:9100"}'
            values: '10+10x7 10+30x4' # 10 20 30 40 50 60 70 80 10 40 70 100 130

      # Unit test for alerting rules.
      alert_rule_test:
          # Unit test 1.
          - eval_time: 10m
            alertname: InstanceDown
            exp_alerts:
                # Alert 1.
                - exp_labels:
                      severity: page
                      instance: localhost:9090
                      job: prometheus
                  exp_annotations:
                      summary: "Instance localhost:9090 down"
                      description: "localhost:9090 of job prometheus has been down for more than 5 minutes."
      # Unit tests for promql expressions.
      promql_expr_test:
          # Unit test 1.
          - expr: go_goroutines > 5
            eval_time: 4m
            exp_samples:
                # Sample 1.
                - labels: 'go_goroutines{job="prometheus",instance="localhost:9090"}'
                  value: 50
                # Sample 2.
                - labels: 'go_goroutines{job="node_exporter",instance="localhost:9100"}'
                  value: 50
```
- the file for testing (alert.yaml)
  ```
    # This is the rules file.

groups:
- name: example
  rules:

  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
        severity: page
    annotations:
        summary: "Instance {{ $labels.instance }} down"
        description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  - alert: AnotherInstanceDown
    expr: up == 0
    for: 10m
    labels:
        severity: page
    annotations:
        summary: "Instance {{ $labels.instance }} down"
        description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."
  ```