# awless-scheduler

The scheduler service is a daemon service that receives templates to be ran and reverted at a later time. 

The service basically get templates, validates and stores them. Then it will check every so often when stored templates need to be executed.

# Usage

### Test

    go test ./... -v

### Run

As an unix sock daemon:

    go build; ./awless-scheduler              # default to scheduler service on unix sock

As an HTTP daemon:

    go build; ./awless-scheduler --http-mode  # default to scheduler service on localhost:8083
    go build; ./awless-scheduler --http-mode --scheduler-hostport 0.0.0.0:9090

Clients use the discovery service to know where the scheduler service is running. By default, the discovery service runs on localhost:8082. To run it on a different port:

    ./awless-scheduler --discovery-hostport localhost:9090

# Usage with the `awless` CLI

The scheduler is mostly used together with the [`awless` CLI](https://github.com/wallix/awless).

From the CLI you can run one-liner, file or remote template and specify the following flags to run your template at a later date:

- `--schedule`: indicates the CLI that this template will be send to the service instead of being scheduled.
- `--run-in`: postpone the execution waiting the `run-in` duration (using [Golang duration notation](https://golang.org/pkg/time/#ParseDuration))
- `--revert-in`: indicates when to revert this template in case it had a succesfull execution

Examples:

    awless create instance name=MyInstance --schedule --run-in 2h --revert-in 4h
    awless create instance name=MyInstance --schedule --revert-in 1d

## Client API

Create a new client giving the discovery URL:

```go
cli, err := client.New("http://127.0.0.1:8082")
```

Behind the scene, the correct client will be instantiated: a UnixSock client or an HTTP client.

Post a template

```go
err := cli.Post(client.Form{
  Region:   "us-west-1",
  RunIn:    "2m",
  RevertIn: "2h",
  Template: txt,
})
```

List tasks

```go
tasks, err := cli.List()
```
