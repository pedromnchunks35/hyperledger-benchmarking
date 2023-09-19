# SCENARIO 1
- The scenario one consists of a single peer and a single organization
- Since the default network already has this components, i need to change the peer core.yaml file for allowing prometheus, the channel config file for only containing one peer and also change the orderer.yaml file for allowing prometheus
- after that config i should reset the peer and the orderer data my deleting the info that they retain in their folders
- In the channel config file i will only leave peer1
- In the core.yaml file, i will open the metrics in the port 9443 and change the metrics to prometheus
- In the orderer, i will open the metrics in the port 8443