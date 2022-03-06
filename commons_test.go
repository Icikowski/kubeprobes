package kubeprobes

import (
	"errors"
	"testing"
	"time"
)

func TestStatusQueryIsAllGreen(t *testing.T) {
	tests := map[string]struct {
		probes         []ProbeFunction
		expectedStatus bool
	}{
		"all green": {
			probes: []ProbeFunction{
				func() error { return nil },
				func() error { time.Sleep(2 * time.Second); return nil },
			},
			expectedStatus: true,
		},
		"some failed": {
			probes: []ProbeFunction{
				func() error { return nil },
				func() error { time.Sleep(2 * time.Second); return errors.New("failed") },
			},
			expectedStatus: false,
		},
		"all failed": {
			probes: []ProbeFunction{
				func() error { return errors.New("failed") },
				func() error { time.Sleep(2 * time.Second); return errors.New("failed") },
			},
			expectedStatus: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sq := newStatusQuery(test.probes)
			if sq.isAllGreen() != test.expectedStatus {
				t.Errorf("expected status %v, got %v", test.expectedStatus, sq.isAllGreen())
			}
		})
	}
}
