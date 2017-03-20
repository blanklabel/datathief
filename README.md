# So what is it?
Datathief WILL be a tool that will connect to open datastorage systems (such as redis and mongo)
pull their version and system information and eventually download all data stored within a database/collection
# Install
    go get github.com/blanklabel/datathief
    
# Configuration
Update the targets.json with the corresponding targets to pull data from

# Results
Results of the run will be placed in the local "datathiefjson.log"

# TODO
 * TODO: Swap to new golang plugin system
 * TODO: Fully concurrent -- connect == current, main app selects between getinfo, dump and connect
 * TODO: Fully completed cmd interface
 * TODO: Cassandra
 * TODO: Elasticsearch
 * TODO: Some datadumper