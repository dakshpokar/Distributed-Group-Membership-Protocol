# MP2-Distributed-Systems

This is a distributed group membership protocol system. The system is implemented using Go, leveraging concurrency and distributed processing across multiple nodes in a network.

## Project Setup for all machines

0. Ensure that you have Go Lang Installed on your machine.
1. Clone this repo on each machine.
2. cd into mp1-distributed-systems folder on all machines.
```bash
cd mp2-distributed-systems
```
3. Start the introducer on machine 1 - 
```go
go run main.go introducer
```
4. Ensure that your log files are present in the same directory on all machines. Example - `/root/`

## How to run Distributed Group Membership Protocol
Now that the introducer is started, you can execute go run main.go and join from any of the machines but before that you have to follow these steps -

1. Create a `.env` file in mp1-distributed-system folder and add the following on all machines-
```env
INTRODUCER_ADDR=127.0.0.1
```
Enter the introducer IP / Hostname. For example - 
```env
INTRODUCER_ADDR=fa24-cs425-9201
```
This has to be performed once on the machine where the distributed disseminator is to be run. If you run distributed 
2. Run main.go on the machine
```bash
go run main.go
```

## Group G92

Rishi Mundada (rishirm3@illinois.edu)
Daksh Pokar (dakshp2@illinois.edu)