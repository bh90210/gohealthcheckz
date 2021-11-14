<p align="center">
  <img width="26%" src="https://user-images.githubusercontent.com/22690219/119139792-1b1b8000-ba4c-11eb-8c88-34c439eada3b.png" />
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/bh90210/healthz.svg)](https://pkg.go.dev/github.com/bh90210/healthz)
[![Go Report Card](https://goreportcard.com/badge/github.com/bh90210/healthz)](https://goreportcard.com/report/github.com/bh90210/healthz)

# healthz

A tiny & extremely simple to use package for Kubernetes liveness/readiness/termination checks.

# Use

_For full examples see the [examples]() folder._

## Init 	

In your `init` or at the begging of your `main` function include:
```go
h := healthz.NewCheck(healthz.OptionsLivePath("live"),
	healthz.OptionsReadyPath("ready"), healthz.OptionsPort("8080"))

go func() {
	if err := h.Start(); err != nil {
		panic(err)
	}
}()
```
_It is a blocking function._

Defining `healthz.NewCheck()` arguments is optional. If nothing's set ie. `healthz.NewCheck("", "", "")` the above default values will be used.

Kubernetes config excerpt:
```yaml
...
        livenessProbe:
          httpGet:
            path: /live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 60
          timeoutSeconds: 5
          failureThreshold: 6
          successThreshold: 1
        startupProbe:
          httpGet:
            path: /live
            port: 8080
          failureThreshold: 15
          periodSeconds: 5
...
```

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

Kubernetes config excerpt:
```yaml
...
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 5
          failureThreshold: 6
          successThreshold: 1
...
```

## Terminating

Kubernetes allows a grace period for any necessary clean up based on your `terminationGracePeriodSeconds` config setting.
```go
if term := hc.Terminating(); term == true {
	// do some clean up
}
```
_It is a blocking function._

Kubernetes config excerpt:
```yaml
...
      terminationGracePeriodSeconds: 5
...
```

## gRPC

The canonical way to achieve health checks with Kubernetes is to include [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe/) along your dockerized gRPC app.

```dockerfile
FROM alpine:3.9
RUN apk add wget
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 
RUN chmod +x /bin/grpc_health_probe  
```

## Init

```go
import healthz "github.com/bh90210/healthz/grpc"

func main() {
hc := healthz.NewCheckGRPC("live", "ready", "5000")

go func() {
	if err := hc.Start(); err != nil {
		// do some error handling
	}
}()
}
```

Kubernetes config excerpt:
```yaml
...
    readinessProbe:
      exec:
        command: ["/bin/grpc_health_probe", "-addr=:5000"]
      initialDelaySeconds: 5
    livenessProbe:
      exec:
        command: ["/bin/grpc_health_probe", "-addr=:5000"]
      initialDelaySeconds: 10
...
```

### Readiness & Terminating

Readiness is section is similar to HTTP. User can utilize `hc.Ready()` & `hc.NotReady()` API to designate the state of the app when probed by Kubernetes.

If user opts to also use the `Terminating` API with gRPC an HTTP server will start.


# Contributing

We are using a feature request workflow. Fork the repo create a new branch ie `fix/http` or `feat/newfeature` and make a PR against `main` branch.

# References

1. https://kubernetes.io/blog/2018/10/01/health-checking-grpc-servers-on-kubernetes/
2. https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
3. https://github.com/grpc-ecosystem/grpc-health-probe/
4. https://pkg.go.dev/google.golang.org/grpc/health/grpc_health_v1?utm_source=godoc
5. https://developers.redhat.com/blog/2020/11/10/you-probably-need-liveness-and-readiness-probes#example_4__putting_it_all_together
6. https://github.com/americanexpress/grpc-k8s-health-check
