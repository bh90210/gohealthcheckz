# healthz

A tiny & extremely simple to use library for Kubernetes liveness/readiness/termination checks.

# Use

## Start

In your `init` or at the begging of your `main` function include:
```go
var healthCheck healthz.Check

go func() {
	if err := healthCheck.Start(); err != nil {
    // Do something with the err
    log.Fatalln(err)
	}
}()
```
_It is a blocking function._

If you would like to specify endpoints and port you can:
```go
healthCheck.Liveness = "live"
healthCheck.Readiness = "ready"
healthCheck.Port = "8080"
```
Defining them is optional. If not set the above default values will be used.

## Ready & NotReady

### Ready

The default state is `NotReady`. When your application is ready to service requests you should set it to:
```go
healthCheck.Ready()
``` 

### NotReady

If your application stops being able to server requests and you wish to let Kubernetes proceed to the appropriate action based on `restartPolicy` (either wait or restart the container) you can achieve that by setting:
```go
healthCheck.NotReady()
```

## Terminating

Kubernetes allows a grace period for any necessary clean up based on your `terminationGracePeriodSeconds` config setting.
```go
if term := healthCheck.Terminating(); term == true {
	// do some clean up
}
```
_It is a blocking function._

## gRPC

WIP

# Contributing

We are using a feature request workflow. Fork the repo create a new branch ie `fix/http` or `feat/newfeature` and make a PR against `main` branch.