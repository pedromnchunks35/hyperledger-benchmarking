# Configuration
- Prometheus is configured by command line flags and a configuration file  
- to specify the configuration file we need to do --config.file
## Generic placeholders
- boolean
- duration (1d,5m,10s)
- filename
- float
- host
- int
- labelname
- labelvalue
- path
- scheme (http or https)
- secret (password)
- string
- size(512MB)
- tmpl_string
## Template
- You can find a template with everything right [here](template.yml)
### Global
- The global configuration specifies parameteres that are valid in all other configuration contexts. Its serves as default configs ofr other sections
  
|Property|Definition|Propertie type/format|default value|
|--|--|--|--|
|scrape_interval|How frequently we shall scrape targets by default|duration|1m|
|scrape_timeout|How long until scrape requests timeout|duration|10s|
|evaluation_interval|How frequently should we evaluate rules|duration|1m|
|external_labels|Labels to apply when series or alerts from external systems come (federation,remote storage, alert manager)|labelname:labelvalue||
|query_log_file|File used for query|string||
|body_size_limit|Limit of size for uncompressed bodys, case it breaks the boundaries, the scrape will fail|size|0, which means no limit|
|sample_limit|Limit of accepted scraped samples from a label, case it breaks, then all scrape is treated as failed|int|0,which means no limit|
|label_limit|Limit of number of labels for a sample.Case it breaks, then all the scrape will be treated as failed|int|0,which means no limit|
|label_name_length_limit|Limit on the label name length,case it breaks all the scrape will be treated as failed|int|0, which means no limit|
|label_value_length_limit|Limit on the value associated with a label length,case it breaks all the scrape will be treated as failed|int|0,which means no limit|
|target_limit|Limit the unique number of targets,case it breaks, all the target will fail|int|0, which means no limit|
|keep_dropped_targets|Limit number of targets dropped by relabeling|int|0, which means no limit|
|rule_files|List of rules and alerts readed by all the matching files|array of path||
|scrape_config_files|Scrape configs readed from files|array of path||
|scrape_configs|List of scrape configs|array of scrape_config||
|alerting|Settings related to the alert manager|array of alert_configs||
|alert_related_configs||array of relabel_config||
|alertmanagers||array of alertmanager_config||
|remote_write|settings related to the remote write feature|array of remote_write||
|remote_read|settings related to the remote read feature|array of remote_read||
|storage|settings related to storage at runtime reloadable|storage_settings||
|tracing|configures exporting traces|||
### scrape_config
- specifies a set of targets and parameteres describing how to scrape them
- targets may be statically configured via the static_configs parameter or dynamically discovered using discovery mechanisms
- relabel_configs allow advanced modifications to any target and its labels before scraping

|Property|Definition|Type/format|Default Value|
|--|--|--|--|
|job_name|The name of the job to scrap metrics|string||
|scrape_interval|overwrite of global_config.scrape_interval|duration|global_config.scrape_interval|
|scrape_timout|overwrite of global_config.scrape_timeout|duration|global_config.scrape_timout|
|scrape_classic_histograms|to scrape a classic histogram as a native|boolean|false|
|metrics_path|The http resource path on which to fetch metrics|path|/metrics|
|honor_labels|How to treat conflits in terms of labels naming, case true it will ignore the conflict, case false it will rename the conflict ones to "exported_(name)"|boolean|false|
|honor_timestamps|if wether prometheus respects the timestamps in the scraped data. case true, the timestamps exposed by the target will be used otherwise they will be ignored|boolean|true|
|scheme|Configures the protocol scheme for requests|scheme|http|
|params|Optional HTTP URL parameteres|string:[string,...]||
|basic_auth|sets the authorization header on every scrape with the configured username and password|username:string<br>password:secret<br>password_file:string||
|authorization|sets the authorization header on every scrape request with the configured credentials|auth_creds||
|authorization.type|sets the authentication type of the request|string|Bearer|
|authorization.credentials|sets the credetials of the request, it is mutually exclusive with credetials_file|string||
|authorization.credetials_file|sets the credentials of the request with the credetials read from a configured file, it is mutually exclusive with credentials|filename||
|oauth2|OAuth 2.0, cannot be used at the same time as basic_auth|oauth2||
|follow_redirects|if scrape requests follow HTTP 3xx redirects|boolean|true|
|enable_http2|if we should enable http2|boolean|true|
|tls_config|tls settings|tls_settings||
|proxy_url|Optional proxy URL|string||
|no_proxy|Comma-separated string that can contain IPS,CIDR notation,domain names,etc, that should be excluded from proxying. IP and domain names can contain port numbers|string||
|proxy_from_enviroment|use proxy URL indicated by enviroment variables (HTTP_PROXY,https_proxy,HTTPs_PROXY,https_proxy,and no_proxy)|Boolean|false|
|proxy_connect_header|Specifies headers to send to proxies during CONNECT requests|string:[(secret),...]||
|azure_sd_configs|List of Azure service discovery configs|array of azure_sd_config||
|consul_sd_configs|List of Consul service discovery configurations|array of consul_sd_config||
|digitalocean_sd_configs|list of digital ocean service discovery configurations|array of digitalocean_sd_config||
|docker_sd_configs|list of docker service discovery configs|array of docker_sd_config||
|dockerswarm_sd_configs|list of docker swarm service discovery configs|array of dockerswarm_sd_config||
|dns_sd_configs|list of dns service discovery configurations|array of dns_sd_config||
|ec2_sd_configs|list of EC2 service discovery configurations|array of ec2_sd_config||
|eureka_sd_configs|list of Eureka service discovery configurations|array of eureka_sd_config||
|file_sd_configs|list of GCE service discovery configurations|array of gce_sd_config||
|gce_sd_configs|list of Hetzmer service discovery configurations|array of hetzner_sd_config||
|http_sd_configs|list of HTTP service discovery configurations||
|ionos_sd_configs|list of IONOS service discovery configs|array of ionos_sd_config||
|kubernetes_sd_configs|List of Kubernetes service discovery configurations|array of kubernetes_sd_config||
|kuma_sd_configs|List of Kuma service discovery configurations|array of kuma_sd_config||
|lightsail_sd_configs|List of Lightsail service discovery configurations|array of lightsail_sd_config||
|linode_sd_configs|List of Linode service discovery configurations|array of linode_sd_config||
|marathon_sd_configs|List of Marathon service discovery configurations|array of marathon_sd_config||
|nerve_sd_configs|List of AirBnB's Nerve service discovery configurations|array of nerve_sd_config||
|nomad_sd_configs|List of Nomad service discovery configurations|array of nomad_sd_config||
|openstack_sd_configs|List of OpenStack service discovery configurations|array of openstack_sd_config||
|ovhcloud_sd_configs|List of OVHcloud service discovery configurations|array of ovhcloud_sd_config||
|puppetdb_sd_configs|List of PuppetDB service discovery configurations|array of puppetdb_sd_config||
|scaleway_sd_configs|List of Scaleway service discovery configurations|array of scaleway_sd_config||
|serverset_sd_configs|List of Zookeeper Serverset service discovery configurations|array of serverset_sd_config||
|triton_sd_configs|List of Triton service discovery configurations|array of triton_sd_config||
|uyuni_sd_configs|List of Uyuni service discovery configurations|array of uyuni_sd_config||
|static_configs|List of labeled statically configured targets for this job|array of static_config||
|relabel_configs|List of target relabel configurations|array of relabel_config||
|metric_relabel_configs|List of metric relabel configurations|array of relabel_config||
|body_size_limit|limit for the size of the response,it will overwrite the one in global_config.body_size_limit|size|0,which means there is no limit|
|sample_limit|Per-scrape limit on the number of scraped samples that will be accepted|int|0,which means there is no limit|
|label_limit|Per-scrape limit on the number of labels that will be accepted for a sample|int|0,which means there is no limit|
|label_name_length_limit|Per-scrape limit on the length of label names that will be accepted for a sample|int|0,which means there is no limit|
|label_value_length_limit|Per-scrape limit on the length of label values that will be accepted for a sample|int|0,which means there is no limit|
|target_limit|Per-scrape config limit on the number of unique targets that will be accepted|int|0,which means there is no limit|
|keep_dropped_targets|Per-job limit on the number of targets dropped by relabeling that will be kept in memory|int|0,which means there is no limit|
|native_histogram_bucket_limit|Limit on the total number of positive and negative buckets allowed in a single native histogram|int|0,which means there is no limit|
### tls_config
|Property|Definition|Type/format|Default Value|
|--|--|--|--|
|ca_file|path for the file that retains a ca identification for the api|filename||
|cert_file|client certificate file path|filename||
|key_file|client key file path|filename||
|server_name|nname of the server|string||
|insecure_skip_verify|disable or not validation of the server certificate|Boolean|false|
|min_version|Accepted min tls version|string||
|max_version|Accepted max tls version|string||
### oauth2
|Property|Definition|Type/format|Default Value|
|--|--|--|--|
|client_id|id of the client|string||
|client_secret|client secret|secret||
|client_secret_file|client secret file to read(must use this prop or the client_secret,not both)|filename||
|scopes|scopes for the token request|array of string||
|token_url|The URL to fetch the token from|string||
|endpoint_params|optional parameteres for the URL|string:string...||
|tls_config|configures the token request tls config|tls_config||
|proxy_url|optional proxy URL|string||
|no_proxy|Comma-separated string that can contain IPs, CIDR notation, domain names that should be excluded from proxying|string||
|proxy_from_enviroment|Use proxy URL indicated by environment variables (HTTP_PROXY, https_proxy, HTTPs_PROXY, https_proxy, and no_proxy)|Boolean|false|
|proxy_connect_header|headers to send to proxies during connect requests|string:[secret,...]||
### static_config
- allows specifying a list of targets and a common label set for them
  
|Property|Definition|Type/format|Default Value|
|--|--|--|--|
|targets| The targets specified by the static config|array of host||
|labels|labels assigned to all metrics scraped from the targts|labelname:labelvalue||
### relabel_config
- Powerful tool to dynamically rewrite the label set of a target before it gets scrapped
- Multiple relabeling steps can be configured
- They are applied to the label set of each target in order of their appearence in the config file

|Property|Definition|Type/format|Default Value|
|--|--|--|--|
|source_labels|source labels select values from existing labels, matching, it is used to replace,keep and drop actions|'['labelname']'||
|separator|the separation of contactenaed source label values|string|;|
|target_label|label to which the result value is written in a replace action<br>It is mandatory for replace actions. Regex capture groups are available.|labelname||
|regex|Regular expression against which the extracted value is matched.|regex|(.*)|
|modulus|Modulus to take of the hash of the source label values|int||
|replacement|Replacement value against which a regex replace is performed if the regular expression matches. Regex capture groups are available.|string|$1|
|action|Action to perform based on regex matching|relabel_action|replace|
<regex> is any valid RE2 regular expression. It is required for the replace, keep, drop, labelmap,labeldrop and labelkeep actions. The regex is anchored on both ends. To un-anchor the regex, use .*<regex>.*.

<relabel_action> determines the relabeling action to take:

- replace: Match regex against the concatenated source_labels. Then, set target_label to replacement, with match group references (${1}, ${2}, ...) in replacement substituted by their value. If regex does not match, no replacement takes place.
- lowercase: Maps the concatenated source_labels to their lower case.
- uppercase: Maps the concatenated source_labels to their upper case.
- keep: Drop targets for which regex does not match the concatenated source_labels.
- drop: Drop targets for which regex matches the concatenated source_labels.
- keepequal: Drop targets for which the concatenated source_labels do not match target_label.
- dropequal: Drop targets for which the concatenated source_labels do match target_label.
- hashmod: Set target_label to the modulus of a hash of the concatenated source_labels.
- labelmap: Match regex against all source label names, not just those specified in source_labels. Then copy the values of the matching labels to label names given by replacement with match group references (${1}, ${2}, ...) in replacement substituted by their value.
- labeldrop: Match regex against all label names. Any label that matches will be removed from the set of labels.
- labelkeep: Match regex against all label names. Any label that does not match will be removed from the set of labels.

Care must be taken with labeldrop and labelkeep to ensure that metrics are still uniquely labeled once the labels are removed.

# [Rules](rules.md)