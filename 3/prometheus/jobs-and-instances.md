# Jobs and instances
- A endpoint you can store data is called an instance
- It means a single process
- A collection of instances of the same purpose, process replicated for scalability or reliability is a job
- Example
  ```
  job: api-server
    -> instance 1: 1.2.3.4:5670
    -> instance 2: 1.2.3.4:5671
    -> instance 3: 1.2.3.4:5672
    -> instance 4: 1.2.3.4:5673
  ```
## Automatically generated labels and time series
- When prometheus saves a target it automatically attaches some labels to the time series which serve to identify the target
- For each isntance saved, prometheus saves a sample in the following time series:
```
up{job="<job-name>", instance="<instance-id>"}: 1 if isntance is healthy or 0 case it failed

scrape_duration_seconds{job="<job-name>", instance="<instance-id>"}: duration of the scrape/save

scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}: the number of samples remaining after metric relabeling was applied

scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}: the number of samples the target exposed

scrape_series_added{job="<job-name>", instance="<instance-id>"}: the approximate number of new series in this scrape
```