# Metric types
- Prometheus client libraries offer four core metric types
## Counter
- Cumulative metric that represents a single monotically increasing number (it means it only increases or stays the same)
- It can only increase or reset to zero
- It can be used for requests served,tasks completed or even errors
## Gauge
- Metric that represents a single numerical value that can arbitrarily go up and down
- Examples of it are temperatures or current memory usage
## Histogram
- Samples observations (usually requests durations or response sizes) and counts them in configurable buckers 
- It groups values from the same metric into several buckets
- We can even count or sum those "categories"
## Summary 
- Much like histograms but more simplier and configurable