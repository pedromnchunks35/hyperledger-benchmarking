# Rate Controllers
- It is a profile that can be either fixed or customizable for specifying at which rate we will send transactions
## Profiles
### Fixed Rate
- The most basic controller
- Default option when no controller is specified
- It will send input transactions at a fixed interval of time specified as TPS(Transactions per secound)
- Example
  ```
  {
    "type": "fixed-rate",
    "opts": {
        "tps": 10 // Transactions per secound that we want
    }
  }
  ```
### Fixed Feedback Rate
- Extension of the fixed Rate
- When the unfinished transactions exceeds times of the defined unfinished transactions for each working, it will stop sending input transactions temporally by sleeping a long period of time
- Example
  ```
  {
  "type": "fixed-feedback-rate",
  "opts": {
      "tps" : 100, // Transactions per secound that we want
      "transactionLoad": 100 // the maximum transaction load on the SUT at which workers will pause sending further transactions (also transactions per secound )
  }
  }
  ```
### Fixed Load
- Profile that aims to produce a given load by changing periodically the given TPS to achieve that load
```
{
  "type": "fixed-load",
  "opts": {
    "transactionLoad": 5, //load we want to achieve
    "startingTps": 100 // the starting tps that later on will be changed
  }
}
```
### Maximum Rate
- This is a profile where we increase the rate until the maximum the SUT can handle until get a certain overload
- Example
```
{
  "type": "maximum-rate",
  "opts": {
    "tps": 100, //the starting transactions per secound
    "step": 5, //the increase of tps for each interval
    "sampleInterval": 20, // the minimum time until tps incrementation
    "includeFailed": true // If we should include failed transactions or not
  }
}
```
### Linear Rate
- It is a controller where we can find a certain load intensity using a interval of TPS
- The TPS is divided by the workers
```
{
  "type": "linear-rate",
  "opts": {
    "startingTps": 25, //lower interval
    "finishingTps": 75 //max interval
    }
}
```
### Composite Rate
- A composition of rates that will switch all over the time
```
{
  "type": "composite-rate",
  "opts": {
    "weights": [2, 1, 2], //Represent the weigth in terms of importance that we give to a certain object, since there are 3 objects inside the array of "rateControllers" there is 3 values inside of it, representing on the same order each object. In this example the unity is time so 2+1+2 = 5 minutes, and which importance means individually one object of the array. Note that we decide if the rate follows a time logic or a number logic
    "rateControllers": [ // array of objects that has different types of controllers
      {
        "type": "fixed-rate",
        "opts": {"tps" : 100}
      },
      {
        "type": "fixed-rate",
        "opts": {"tps" : 300}
      },
       {
        "type": "fixed-rate",
        "opts": {"tps" : 200}
      }
    ],
    "logChange": true // Show a log in case we switch controller
  }
}
```
### Zero Rate
- This controlled is meant to execute each rate controller for a specific period of time (this only works if the benchmarking is time based). The logic is equal to the composed rate, but this is just for time. The options are almost the same, de difference is that the weigths is only for time
```
{
  "type": "composite-rate",
  "opts": {
    "weights": [30, 10, 10, 30],
    "rateControllers": [
      {
        "type": "fixed-rate",
        "opts": {"tps" : 100}
      },
      {
        "type": "fixed-rate",
        "opts": {"tps" : 500}
      },
      {
        "type": "zero-rate",
        "opts": { }
      },
      {
        "type": "fixed-rate",
        "opts": {"tps" : 100}
      }
    ],
    "logChange": true
  }
}
```
### Record Rate
- This a rate that enables us to store transactions in a file
```
{
  "type": "record-rate",
  "opts": {
    "rateController": {
      "type": "fixed-rate",
      "opts": {"tps" : 100}
    },
    "pathTemplate": "../tx_records_client<C>_round<R>.txt",// <C> stands for the placeholder of the worker and <R> stands for the placeholder of the current round
    "outputFormat": "TEXT", // Type of file, it can be either text or binary  
    "logEnd": true //IF LOGS MUST BE ON THE FILE
  }
}
```
### Replay Rate
- This is a rate for reproducing the same workload behavior that we got in previous saved in file rates
```
{
  "type": "replay-rate",
  "opts": {
    "pathTemplate": "../tx_records_client<C>.txt", // Path to a existing file to replicate behavior
    "inputFormat": "TEXT",
    "logWarnings": true,
    "defaultSleepTime": 50 // THe sleep time between transaction, default is 20 ms
    }
}
```
### Adding custom controllers
- In order to add custom constrollers we need to implementent interfaces in a javascript file implementing the desired behavior