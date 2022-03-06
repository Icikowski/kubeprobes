package kubeprobes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func getStatusFromEndpoint(t *testing.T, client *http.Client, endpoint string) int {
	t.Helper()
	resp, err := client.Get(endpoint)
	if err != nil {
		t.Errorf("error getting status from endpoint: %s", err)
	}
	return resp.StatusCode
}

func TestKubeprobes(t *testing.T) {
	live, ready := NewStatefulProbe(), NewStatefulProbe()

	tests := map[string]struct {
		livenessProbeTransformation  func(*testing.T, *statefulProbe)
		readinessProbeTransformation func(*testing.T, *statefulProbe)
		expectedLiveStatus           int
		expectedReadyStatus          int
	}{
		"not live": {
			livenessProbeTransformation:  markAsDown,
			readinessProbeTransformation: markAsDown,
			expectedLiveStatus:           http.StatusServiceUnavailable,
			expectedReadyStatus:          http.StatusServiceUnavailable,
		},
		"live but not ready": {
			livenessProbeTransformation:  markAsUp,
			readinessProbeTransformation: markAsDown,
			expectedLiveStatus:           http.StatusOK,
			expectedReadyStatus:          http.StatusServiceUnavailable,
		},
		"live and ready": {
			livenessProbeTransformation:  markAsUp,
			readinessProbeTransformation: markAsUp,
			expectedLiveStatus:           http.StatusOK,
			expectedReadyStatus:          http.StatusOK,
		},
		"ready but not live - should never happen": {
			livenessProbeTransformation:  markAsDown,
			readinessProbeTransformation: markAsUp,
			expectedLiveStatus:           http.StatusServiceUnavailable,
			expectedReadyStatus:          http.StatusServiceUnavailable,
		},
	}

	kp := NewKubeprobes(
		WithLivenessProbes(live.GetProbeFunction()),
		WithReadinessProbes(ready.GetProbeFunction()),
	)

	srv := httptest.NewServer(kp)
	defer srv.Close()
	client := srv.Client()

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			test.livenessProbeTransformation(t, live)
			test.readinessProbeTransformation(t, ready)

			liveStatus := getStatusFromEndpoint(t, client, srv.URL+"/live")
			readyStatus := getStatusFromEndpoint(t, client, srv.URL+"/ready")
			otherStatus := getStatusFromEndpoint(t, client, srv.URL+"/something")

			if liveStatus != test.expectedLiveStatus {
				t.Errorf("expected live status %d, got %d", test.expectedLiveStatus, liveStatus)
			}
			if readyStatus != test.expectedReadyStatus {
				t.Errorf("expected ready status %d, got %d", test.expectedReadyStatus, readyStatus)
			}
			if otherStatus != http.StatusNotFound {
				t.Errorf("expected 404 status, got %d", otherStatus)
			}
		})
	}
}
