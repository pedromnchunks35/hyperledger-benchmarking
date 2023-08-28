# Logging control
- In calliper you can customize your logging style (in the runtime config)
- Config the logging targets (in the runtime config)
- Create your own loggers
## Customizing the logging style
```
caliper:
  logging:
    template: '%timestamp% %level% [%label%] [%module%] %message% (%metadata%)'
```
- This is setted in the caliper config file.. you can provide the config flag as a flag like so:
  ```
  npx caliper benchmark run -c path/to/caliper-config.yaml
  ```
|Placeholder|Required format|Description|
|--|--|--|
|%timestamp%|timestamp|it will be replaced with the timestamp of the log message|
|%level%|-|Will be replaced with the severity level (info,warn,error)|
|%label%|label|It will be replaced with the configuration label of the process|
|%module%|-|It will be replaced with the module name|
|%message%|-|Will be replaced with the message|
|%metadata%|-|Will be replaced with the string representation of additional logging arguments|
- You can also override this template by changing the caliper-logging-template setting key, for example, from the command line: --caliper-logging-template="%time%: %message%"
- Color can be changed under "caliper.logging.formats"
- Example is to use something like this: --caliper-logging-formats-colorize-colors-info=blue
- The messages can be outputed as json using this: "--caliper-logging-formats-json="{space:0}"". This space property is from the JSON.stringify
- Adding padding to the logs "--caliper-logging-formats-pad=true"
- Aligning the logs: --caliper-logging-formats-align=true
- We can retrieve the type of the message using "%attribute%"
- Customizing the level info of the log example: "--caliper-logging-formats-attributeformat-level="LEVEL[%attribute%]""
## Creating own loggers
```
const logger = require('@hyperledger/caliper-core').CaliperUtils.getLogger('my-module');

// ...

logger.debug('My custom debug message', metadataObject1, metadataObject2);
```
## Examples
```
template: '%timestamp%%level%%label%%module%%message%%metadata%'
formats:
    timestamp: 'YYYY.MM.DD-HH:mm:ss.SSS'
    label: caliper
    json: false
    pad: true
    align: false
    attributeformat:
        level: ' %attribute%'
        label: ' [%attribute%]'
        module: ' [%attribute%] '
        metadata: ' (%attribute%)'
    colorize:
        all: true
        colors:
            info: green
            error: red
            warn: yellow
            debug: grey
```
```
caliper:
  logging:
    # no need for timestamp and label
    template: '%level% [%module%]: %message% %meta%'
    formats:
      # color codes look ugly in log files
      colorize: false
      # don't need these, since won't appear in the template
      label: false
      timestamp: false
    targets:
      file:
        options:
          # bump the log level from debug to warn, only log the critical stuff in this file
          level: warn
          filename: 'critical.log'
      rotatingfile:
        target: daily-rotate-file
        enabled: true
        options:
          level: debug
          datePattern: 'YYYY-MM-DD-HH'
          zippedArchive: true
          filename: 'debug-%DATE%.log'
          options:
            flags: a
            mode: 0666
```