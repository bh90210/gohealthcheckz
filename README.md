# healthz

A tiny & extremely simple to use library for Kubernetes liveness/readiness/termination checks.

# Use

## Start

In your `init` or at the begging of your `main` function include:
```go
	go func() {
		if err := healthz.Start(); err != nil {
      // Do something with the err
			log.Fatalln(err)
		}
	}()
```
It is a blocking function so use it
accordingly. 

## Ready & NotReady

### Ready

The default state is `NotReady`. When your application is ready to service requests you should set it to:
```go
healthz.Ready()
``` 

### NotReady

If your application stops being able to server requests and you wish to let Kubernetes proceed to the appropriate action based on `restartPolicy` you can achieve that by setting:
```go
healthz.NotReady()
```

## Terminating

Kubernetes allows a grace period for any necessary clean up based on your `terminationGracePeriodSeconds` config setting.
```go
	if term := healthz.Terminating(); term == true {
		// do some clean up
	}
```
It is a blocking function.

## gRPC

WIP

# Contributing

We are using a feature request workflow. Fork the repo create a new branch ie `fix/http` or `feat/newfeature` and make a PR against `main` branch.