package svc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//nolint:dupl
func TestAlive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		givenError   error
		expectedCode int
	}{
		{
			name: "should return status ok when no error",

			givenError: nil,

			expectedCode: 200,
		},
		{
			name: "should return status not available when an error",

			givenError: fmt.Errorf("internal error, stop taking traffic until ready again"),

			expectedCode: 503,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dummyWorker := &WorkerMock{
				RunFunc: func() error {
					return nil
				},
				TerminateFunc: func() error {
					return nil
				},
				InitFunc: func(*zap.Logger) error { return nil },
				HealthyFunc: func() error {
					return tc.givenError
				},
			}

			s, err := New("dummy-service", "v0.0.0", WithHealthz(), WithHTTPServer("9090"))
			require.NoError(t, err)

			s.AddWorker("dummy-worker", dummyWorker)

			go s.Run()

			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			rec := httptest.NewRecorder()
			s.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

//nolint:dupl
func TestHealthy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		givenError   error
		expectedCode int
	}{
		{
			name: "should return status ok when no error",

			givenError: nil,

			expectedCode: 200,
		},
		{
			name: "should return status not available when an error",

			givenError: fmt.Errorf("internal error, restart container"),

			expectedCode: 503,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dummyWorker := &WorkerMock{
				RunFunc: func() error {
					return nil
				},
				TerminateFunc: func() error {
					return nil
				},
				InitFunc: func(*zap.Logger) error { return nil },
				AliveFunc: func() error {
					return tc.givenError
				},
			}

			s, err := New("dummy-service", "v0.0.0", WithHealthz(), WithHTTPServer("9090"))
			require.NoError(t, err)

			s.AddWorker("dummy-worker", dummyWorker)

			go s.Run()

			req := httptest.NewRequest(http.MethodGet, "/live", nil)
			rec := httptest.NewRecorder()
			s.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
