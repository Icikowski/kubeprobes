package kubeprobes

import (
	"net/http"
)

type kubeprobes struct {
	livenessProbes  []ProbeFunction
	readinessProbes []ProbeFunction
}

// ServeHTTP implements http.Handler interface
func (kp *kubeprobes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/live":
		sq := newStatusQuery(kp.livenessProbes)
		if sq.isAllGreen() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	case "/ready":
		sq := newStatusQuery(append(kp.livenessProbes, kp.readinessProbes...))
		if sq.isAllGreen() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

type option func(*kubeprobes)

// New returns a new instance of a Kubernetes probes
func New(options ...option) *kubeprobes {
	kp := &kubeprobes{
		livenessProbes:  []ProbeFunction{},
		readinessProbes: []ProbeFunction{},
	}

	for _, option := range options {
		option(kp)
	}

	return kp
}

// WithLivenessProbes adds given liveness probes to the set of probes
func WithLivenessProbes(probes ...ProbeFunction) option {
	return func(kp *kubeprobes) {
		kp.livenessProbes = append(kp.livenessProbes, probes...)
	}
}

// WithReadinessProbes adds given readiness probes to the set of probes
func WithReadinessProbes(probes ...ProbeFunction) option {
	return func(kp *kubeprobes) {
		kp.readinessProbes = append(kp.readinessProbes, probes...)
	}
}
