# Messengers
- Caliper uses a orchestrator to control workers that iteract with the SUT in order to perform a benchmark. 
- Messages are passed between the orchestrator and all workers to keep workers sync.
- User may specify the messaging procol that will be used by caliper in order to make communications easier between the orchestrator and the worker.
## Messengers protocols
- The messaging protocol to be used by caliper is established in the runtime configuration file.
- Permitted messengers are:
    - **process**, which is the default and is based in native NodeJS **process** based communications. This is only valid when the workers are local workers
    - **mqtt**, which is a messenger that uses **MQTT** (publisher subscriver pattern tool) to make communication between workers and the orchestrator easier. It is valid for both local and distributed workers, and assumes the existence of a MQTT broker service that may be used such as mosquito.
- In order to change into a mqtt message protocol you use this:
  ```
      worker:
        communication:
            method: mqtt
            address: mqtt://localhost:1883
  ```
- You can also pass it in the command line like so "--caliper-worker-communication-method mqtt --caliper-worker-communication-address mqtt://localhost:1883"