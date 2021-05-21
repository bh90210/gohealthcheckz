<p align="center">
  <img width="26%" src="https://user-images.githubusercontent.com/22690219/119138540-a8f66b80-ba4a-11eb-8fc9-273a4328db1d.png" />
</p>
 
[![Go Report Card](https://goreportcard.com/badge/github.com/bh90210/healthz)](https://goreportcard.com/report/github.com/bh90210/healthz)
[![Build Status](https://drone.euoe.dev/api/badges/bh90210/healthz/status.svg)](https://drone.euoe.dev/bh90210/healthz)
[![codecov](https://codecov.io/gh/bh90210/healthz/branch/main/graph/badge.svg?token=9PSK4W6VJ9)](https://codecov.io/gh/bh90210/healthz)
[![DeepSource](https://deepsource.io/gh/bh90210/healthz.svg/?label=active+issues&show_trend=true)](https://deepsource.io/gh/bh90210/healthz/?ref=repository-badge)

# healthz

A tiny & extremely simple to use package for Kubernetes liveness/readiness/termination checks.

# Use

## Start

In your `init` or at the begging of your `main` function include:
```go
hc := healthz.NewCheck("live", "ready", "8080")

go func() {
	if err := hc.Start(); err != nil {
		// do some error handling
	}
}()
```
_It is a blocking function._

Defining `healthz.NewCheck()` arguments is optional. If nothing's set ie. `healthz.NewCheck("", "", "")` the above default values will be used.

## Ready & NotReady

### Ready

The default state is `NotReady`. When your application is ready to service requests you should set it to:
```go
hc.Ready()
``` 

### NotReady

If your application stops being able to server requests and you wish to let Kubernetes proceed to the appropriate action based on `restartPolicy` (either wait or restart the container) you can achieve that by setting:
```go
hc.NotReady()
```

## Terminating

Kubernetes allows a grace period for any necessary clean up based on your `terminationGracePeriodSeconds` config setting.
```go
if term := hc.Terminating(); term == true {
	// do some clean up
}
```
_It is a blocking function._

## gRPC

WIP

# Contributing

We are using a feature request workflow. Fork the repo create a new branch ie `fix/http` or `feat/newfeature` and make a PR against `main` branch.
