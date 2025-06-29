# Golang API

<!--toc:start-->

- [Golang API](#golang-api)
  - [Modules](#modules)
    - [Creating HTTP Server](#creating-http-server)
    - [Parsing Command Line Arguments](#parsing-command-line-arguments)
    <!--toc:end-->

We are going to build a Golang Restful API from scratch.

## Modules

### Creating HTTP Server

We are using `net/http` package to create a simple HTTP server.
And we have setup the internal package `internal/app` to hold some
of the application wide configurations.

And listen the server on port `8080` and added health check
route `/health` to check the server status.

### Parsing Command Line Arguments

We are using `flag` package to parse command line arguments.
To dynamically change the port of the server, we have added
port flag to the command line arguments.

```go
var port int
flag.IntVar(&port, "port", 8080, "This is the default port on which the server will run")
flag.Parse()
```

After that we can pass the port number while running the server.

```bash
go run main.go --port=9090
```

### Chi Router

We can use `chi` router to handle the routing in a more structured way.

```go
r := chi.NewRouter()
```
