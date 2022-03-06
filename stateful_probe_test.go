package kubeprobes

import "testing"

var (
	markAsDown func(*testing.T, *StatefulProbe) = func(t *testing.T, sp *StatefulProbe) {
		t.Helper()
		sp.MarkAsDown()
	}
	markAsUp func(*testing.T, *StatefulProbe) = func(t *testing.T, sp *StatefulProbe) {
		t.Helper()
		sp.MarkAsUp()
	}
)

func TestStatefulProbe(t *testing.T) {
	tests := map[string]struct {
		probeTransformation func(*testing.T, *StatefulProbe)
		expectedError       bool
	}{
		"mark as up": {
			probeTransformation: markAsUp,
			expectedError:       false,
		},
		"mark as down": {
			probeTransformation: markAsDown,
			expectedError:       true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			sp := NewStatefulProbe()
			test.probeTransformation(t, sp)
			probeFunc := sp.GetProbeFunction()
			if (probeFunc() != nil) != test.expectedError {
				t.Error("result not as expected")
			}
		})
	}
}
