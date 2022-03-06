# kubeprobes

Simple and effective package for implementing [Kubernetes liveness and readiness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)' handler.

![Go version](https://img.shields.io/github/go-mod/go-version/Icikowski/kubeprobes)
[![Go Report Card](https://goreportcard.com/badge/github.com/Icikowski/kubeprobes)](https://goreportcard.com/report/github.com/Icikowski/kubeprobes)
[![Go Reference](https://pkg.go.dev/badge/github.com/Icikowski/kubeprobes.svg)](https://pkg.go.dev/github.com/Icikowski/kubeprobes)
[![Codecov](https://img.shields.io/codecov/c/gh/Icikowski/kubeprobes?token=85PO16238X)](https://codecov.io/gh/Icikowski/kubeprobes)
![License](https://img.shields.io/github/license/Icikowski/kubeprobes)

## Installation

```bash
go get -u github.com/Icikowski/kubeprobes
```

## Usage

The package provides `kubeprobes.NewKubeprobes` function which returns a probes handler compliant with `http.Handler` interface. 

The handler serves two endpoints, which are used to implement liveness and readiness probes by returning either `200` (healthy) or `503` (unhealthy) status: 

- `/live` - endpoint for liveness probe;
- `/ready` - endpoint for readiness probe.

Accessing any other endpoint will return `404` status. In order to provide maximum performance, no body is ever returned.

The `kubeprobes.NewKubeprobes` function accepts following options-applying functions as arguments:

- `kubeprobes.WithLivenessProbes(/* ... */)` - adds particular [probes](#probes) to the list of liveness probes;
- `kubeprobes.WithReadinessProbes(/* ... */)` - adds particular [probes](#probes) to the list of readiness probes.

## Probes

In order to determine the state of particular element of application, probes need to be implemented either by creating [status determining function](#probe-functions) or by using simple and thread-safe [stateful probes](#stateful-probes). 

### Probe functions

Probe functions (objects of type `ProbeFunction`) are functions that performs user defined logic in order to determine whether the probe should be marked as healthy or not. Those functions should take no arguments and return error (if no error is returned, the probe is considered to be healthy; if error is returned, the probe is considered to be unhealthy).

```go
someProbe := func() error {
    // Some logic here
    if somethingIsWrong {
        return errors.New("something is wrong")
    }
    return nil
}

someOtherProbe := func() error {
    // Always healthy
    return nil
} 

// Use functions in probes handler
kp := kubeprobes.NewKubeprobes(
    kubeprobes.WithLivenessProbes(someOtherProbe),
    kubeprobes.WithReadinessProbes(someProbe),
)
```

### Stateful probes

Stateful probes (objects of type `StatefulProbe`) are objects that can be marked either as "up" (healthy) or "down" (unhealthy) and provide a `ProbeFunction` for easy integration. Those objects utilize `sync.Mutex` mechanism to provide thread-safety.

```go
// Unhealthy by default
someProbe := kubeprobes.NewStatefulProbe()
someOtherProbe := kubeprobes.NewStatefulProbe()

// Use it in probes handler
kp := kubeprobes.NewKubeprobes(
    kubeprobes.WithLivenessProbes(someProbe.GetProbeFunction()),
    kubeprobes.WithReadinessProbes(someOtherProbe.GetProbeFunction()),
)
```

## Example usage

```go
// Create stateful probes
live := kubeprobes.NewStatefulProbe() 
ready := kubeprobes.NewStatefulProbe()

// Prepare handler
kp := kubeprobes.NewKubeprobes(
    kubeprobes.WithLivenessProbes(live.GetProbeFunction()),
    kubeprobes.WithReadinessProbes(ready.GetProbeFunction()),
)

// Start the probes server
probes := &http.Server{
    Addr:    ":8080",
    Handler: kp,
}
go probes.ListenAndServe()

// Mark probes as healthy
live.MarkAsUp()
ready.MarkAsUp()
```
