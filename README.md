# Ormen 
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/ashish493/ormen/issues)
 [![Go Report Card](https://goreportcard.com/badge/github.com/ashish493/ormen)](https://goreportcard.com/report/github.com/ashish493/ormen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<p align="center">
  <img width="490" height="450" src="https://github.com/ashish493/ormen/assets/44671044/3adcf224-922c-41a8-ba4c-788830fcb3ce">
</p>

## About the Project 
Ormen is a toy Orchestrator written in go based on Vikings theme. The name Ormen is derived from Ormen Lange, which was one of the most famous of the Viking longships. It was also known as the The Long Serpent, and was the largest and most powerful longship of its day. Longships were the epitome of Scandinavian naval power at the time and were used by the Norse in warfare. 

### Relation between ships and orchestrator

There is an analogical similarity between a ship and an orchestrator. Both are responsible for coordinating various components, adapting to changing 
conditions, optimizing performance, and ultimately ensuring the successful achievement of their respective goals. 

- Ships are designed for transportation where it navigates through complex environments, overcoming challenges like weather conditions and varying water depths whereas the Orchestrators coordinate and manage the flow of tasks and services, navigating through the intricacies of software deployment, resource allocation, and overall system efficiency.

- Ships optimize routes, manage fuel consumption, and maintain the vessel to ensure it operates at peak performance. Orchestrators optimize IT workflows by automating tasks, managing resources efficiently, and ensuring that applications run smoothly, minimizing downtime and maximizing performance.

- Ships adapt to changing conditions at sea, adjusting their course based on factors like weather, currents, and potential obstacles and Orchestrators adapt to changes in the IT environment, dynamically allocating resources, adjusting workflows, and responding to changes in demand or system conditions.

Below is the Analogical comparison between ships and orchestrator.

<p align="center">
  <img width="500" height="370" src="https://github.com/ashish493/ormen/assets/44671044/a6cc7439-6b31-4527-b637-3ecc8b1d1cd3">
</p>

### Technical Details

This project is completely written in go. 

- Bolt Db is used for persistent storage of tasks. 
- Cobra to manage CLI commands
- chi to manage api routings 
- goprocinfo to monitor the cpu usage
- Docker SDK to manage containers  

## Installation 

```
go get github.com/ashish493/ormen
```

## Requirements
Since we are orchestrating the Docker containers, we need a running Docker Daemon before running this application. 

## Usage

1. Use `go run main.go` to get the list of available commands.

```
$   go run main.go 

    A CLI based toy orchestrator used to manage docker containers written in go.

    Usage:
    ormen [command]

    Available Commands:
    deck        command to operate a manager node.
    help        Help about any command
    mast        Mast command to list nodes.
    run         Run a new task.
    sailor      Worker command to operate a Cube worker node.
    status      Status command to list tasks.
    stop        Stop a running task.

    Flags:
    -h, --help   help for ormen

    Use "ormen [command] --help" for more information about a command.
```

2. Use `go run main.go sailor -p 8081` to setup a sailor node at a specific port number. The default port number of sailor is 5556.

```
$   go run main.go sailor -p 8081
    2023/11/10 18:26:02 Starting worker.
    2023/11/10 18:26:02 Starting worker API on http://0.0.0.0:8081
    api started2023/11/10 18:26:02 No tasks to process currently.
    2023/11/10 18:26:02 Sleeping for 10 seconds.
    2023/11/10 18:26:02 Collecting stats
    2023/11/10 18:26:02 Checking status of tasks
    2023/11/10 18:26:02 Task updates completed
    2023/11/10 18:26:02 Sleeping for 15 seconds
    2023/11/10 18:26:12 No tasks to process currently.
    2023/11/10 18:26:12 Sleeping for 10 seconds.
```

3. Use `go run main.go deck -w 'localhost:8080,localhost:8081,localhost:8082'` to setup a deck for the respective sailor nodes. The default port number of deck is 5555.

```
$ go run main.go deck -w 'localhost:8080,localhost:8081,localhost:8082'
2023/11/10 18:37:47 Starting manager.
2023/11/10 18:37:47 Starting manager API on http://0.0.0.0:5555
2023/11/10 18:37:47 Checking for task updates from workers
2023/11/10 18:37:47 Checking worker localhost:8080 for task updates
2023/11/10 18:37:47 Processing any tasks in the queue
2023/11/10 18:37:47 No work in the queue
2023/11/10 18:37:47 Sleeping for 10 seconds
2023/11/10 18:37:47 Performing task health check
2023/11/10 18:37:47 Task health checks completed
2023/11/10 18:37:47 Sleeping for 60 seconds
2023/11/10 18:37:47 Collecting stats for node localhost:8080
2023/11/10 18:37:47 Checking worker localhost:8081 for task updates
2023/11/10 18:37:47 Collecting stats for node localhost:8081
2023/11/10 18:37:47 Collecting stats for node localhost:8082
``` 


4. Use `go run main.go stop bb1d59ef-9fc1-4e4b-a44d-db571eeed203` to stop 

5.  Use `go run main.go mast` to get the list of nodes.  
```
$ go run main.go mast 
NAME               MEMORY (MiB)     DISK (GiB)     ROLE       TASKS     
localhost:8080     7967             1081           worker     0         
localhost:8081     7967             1081           worker     0         
localhost:8082     7967             1081           worker     0
```


4. Use Curl operations on deck's server to perform CRUD Operation on containers and also to view nodes.

```
$ curl -X POST localhost:5555/tasks -d @task1.json

{"ID":"bb1d59ef-9fc1-4e4b-a44d-db571eeed203","ContainerID":"","Name":"test-chapter-9.1","State":1,"Image":"timboring/echo-server:latest","Cpu":0,"Memory":0,"Disk":0,"ExposedPorts":{"7777/tcp":{}},"HostPorts":null,"PortBindings":{"7777/tcp":"7777"},"RestartPolicy":"","StartTime":"0001-01-01T00:00:00Z","FinishTime":"0001-01-01T00:00:00Z","HealthCheck":"/health","RestartCount":0}

$ curl -X DELETE localhost:5555/tasks/bb1d59ef-9fc1-4e4b-a44d-db571eeed203

$ curl localhost:5555/nodes|jq
[
  {
    "Name": "localhost:8080 ",
    "Ip": "",
    "Api": "http://localhost:8080 ",
    "Memory": 32793076,
    "MemoryAllocated": 0,
    "Disk": 20411170816,
    "DiskAllocated": 0,
    "Stats": {
      "MemStats": {...},
      "DiskStats": {...},
      "CpuStats": {...},
      "LoadStats": {...},
      "TaskCount": 0
    },
    "Role": "worker",
    "TaskCount": 0
  }
]

```

You can refer the [task1.json](https://github.com/ashish493/ormen/blob/main/task1.json) to create or update the Sails. 


## Acknowledgement

I took inspiration of this project from the book [Build an Orchestrator in Go (From Scratch)](https://www.simonandschuster.com/books/Build-an-Orchestrator-in-Go-(From-Scratch)/Tim-Boring/From-Scratch/9781617299759). I would like to thank [Tim Boring](https://www.simonandschuster.com/authors/Tim-Boring/191341900) for writing this book in a very detailed manner, which really helped me in building this project. 
