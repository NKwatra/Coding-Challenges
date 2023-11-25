#### lb

- Custom impementation of application load balancer based on round robin scheduling algorithm.
- Provides `-d`, `-c` and `-p` flags for configuring the web server domains, the health check interval (_10s default_) and health check path respectively(_root `/` default_).
- Implements a health check timeout of `2s`. Supports automatic removal and readdition of servers failing health checks.
- Listens on the port 8080 by default.

##### Test Examples

###### To Build & Install

```bash
cd lb
go build && go install
```

###### Execute Post Installation

###### Starting Load Balancer

```bash
    lb -d "http://localhost:3000,http://localhost:3001"
```

```bash
    lb -d "http://localhost:3000,http://localhost:3001,http://localhost:3002" -c 15 -p "/check"
```

###### Testing Servers

> **NOTE**: Please ensure web servers are running on input domains else requests might fail.

```bash
    seq 1 50 | xargs -Iname -P5 curl "http://localhost:8080"
```

##### Concepts Used

- Concurrency
  - Channels
  - Mutex
- Http Servers
- Http Headers Parsing
- Closures
- Http Clients
- Background threads
