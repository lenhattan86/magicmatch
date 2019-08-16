# MagicMatch
MagicMatch is a independent service that does the load based task-to-host matching for Apache Aurora. Traditionally, Aurora does the allocation based on resource requests (CPU, memory, etc.). However, the resource requests are not the same as the resource usage. For example, a task may use more resource than its request. Hence, a busy host may have to accomodate more tasks. Meaning, some tasks unfairly suffer from performance degradation.

# How to build MagicMatch 
## Prerequisites
- Golang
- A modified version of Gorealis https://github.paypal.com/nhatle/gorealis. This version does the Thrift calls for MagicMatch. Gorealis is a library so you don't need to install it.
- Apache Aurora & Mesos
## build commands
- To build MagicMatch for your local OS, you just run "go build -o magicmatch main.go". This will generate an executable file magicmatch.
- To build for other OSes, refer ubuntu_build.sh as an example for Ubuntu.

# How to run MagicMatch 
- Simply run "./magicmatch" it will create a web server on port 8080. You can indicate the port number using -port=9090. 
- More configuration parameters can be found in metrics/configuration.go.