package kubeprobes

import (
	"errors"
	"sync"
)

var errProbeDown = errors.New("DOWN")

type statefulProbe struct {
	status bool
	mux    sync.Mutex
}

// NewStatefulProbe returns a new instance of a stateful probe
// which can be either marked as "up" (healthy) or "down" (unhealthy).
// The probe is initially marked as "down".
func NewStatefulProbe() *statefulProbe {
	return &statefulProbe{
		status: false,
		mux:    sync.Mutex{},
	}
}

// MarkAsUp marks the probe as healthy
func (sp *statefulProbe) MarkAsUp() {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	sp.status = true
}

// MarkAsDown marks the probe as unhealthy
func (sp *statefulProbe) MarkAsDown() {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	sp.status = false
}

// GetProbeFunction returns a function that can be used to check
// whether the probe is healthy or not.
func (sp *statefulProbe) GetProbeFunction() ProbeFunction {
	return func() error {
		sp.mux.Lock()
		defer sp.mux.Unlock()
		if sp.status {
			return nil
		}
		return errProbeDown
	}
}
