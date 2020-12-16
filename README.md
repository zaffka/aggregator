# Simple message queue
It is a test assignment I've got from a potential employer.  
Our negotiations had finished before I started the realization, but I found the assignment interesting and spent some time playing with it. 

I've found some of the preconditions incorrect, so the realization does not 100% fit the conditions.

I've simplified the queue's topics initialization process to avoid mutexes usage. Topics are creating during the queue instantiation, so I don't need to dynamically make them when the aggregators or generators handle the queue.

More of that, the assignment doesn't have enough unit tests, storage part remained unrealized - I didn't want to spend too much time.

But, if you are interested - you can check all the stuff below.

# How to run
Just say `make run`.  
The command will make Docker image and runs the container using a compose file.  
Of course, you need docker and docker-compose installed.  

Use `docker logs [container]` to read logs of the running app, `make down` to stop the container with the app.

# Project structure
```
.
├── config                  # Configuration file parser package
│   ├── bad_conf_file.json
│   ├── config.go
│   ├── config_test.go
│   └── ok_conf_file.json
├── domain                  # App's shared types and interfaces
│   ├── interfaces.go
│   └── types.go
├── exec                    # Package to execute queue's publishers\subscribers
│   ├── exec.go
│   └── internal            # Set of executors
│       ├── aggregator.go
│       ├── aggregators.go
│       ├── generator.go
│       ├── generators.go
│       └── helpers.go
├── log                     # Uber's zap logger init
│   └── log.go
├── queue                   # Basic queue
│   └── queue.go
├── shutdown                # Shutdown control package
│   └── shutdown.go
└── storage                 # Storage package UNREALIZED
    └── storage.go
```

# Architecture
![scheme](scheme.jpg)
# Input config

```
{
    "generators": [
        {
            "timeout_s":                30,         // generator stop timeout in seconds
            "send_period_s":            1,          // data send period in seconds
            "data_sources": [
                {
                    "id":               "data_1",   // data source identificator (dataId)
                    "init_value":       50,         // start value
                    "max_change_step":  5,          // maximum change step for previous value
                },
                ...
            ],
        },
        ...
    ],
    "aggregators": [
        {
            "sub_ids":                  ["data_1],  // array of dataIds that aggregator subscribes to
            "aggregate_period_s":        10          // period in seconds to collect and process data
        }
    ],
    "queue": {
        "size":                         50,         // message queue size
    },
    "storage_type":                     0           // 0 - console, 1 - file with some name
}
```

# Description

## Generator

Starts and produces data for specified period of time (`timeout_s`).
Every generator might produce several data sources (dataIds).
For every data source, value changes in time starting with initial point (`init_value`) and increases randomly with max step (`max_change_step`).

## Message queue

Basic FIFO buffer with single input and multiple output.
Queue should be implemented as pub / sub broker.
Once all generators are stopped, queue must be notified about that event in order to notify all output listeners.

## Aggregator

Starts and aggregates data while input channel is "alive".
Aggregation processing:

- collect all incoming data for specified period of time (`agregate_period_s`);
- calculate average value for every dataId;
- send average values to Storage;
- start new aggregation iteration;
  Sending data to Storage should be asynchronous.

## Data format

{
"id": "data_1",
"value": 50
}

## Abstract storage

Unified interface to store integer value with corresponding name (id).
Every value added to the Storage is a new line (new entry), so no aggregation here.
Storage type is specified in config file.

# Requirements

- only Golang standard libraries. App could be written with the help of Cobra framework. No queue brokers;
- application should react to SIGINT UNIX signal and stop all workers gracefully;
- application should be packed into docker container;
- application should not require any preparation steps, but run docker container, with some parameters if needed;
- application should have short Readme with application running instructions.
