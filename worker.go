package svc

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

// Worker defines a SVC worker.
type Worker interface {
	Init(*slog.Logger) error
	Run() error
	Terminate() error
}

// Aliver defines a worker that can report his livez status.
type Aliver interface {
	Alive() error
}

// Healther defines a worker that can report his healthz status.
type Healther interface {
	Healthy() error
}

// Gatherer is a place for workers to return a prometheus.Gatherer
// for SVC to serve on the metrics endpoint.
type Gatherer interface {
	Gatherer() prometheus.Gatherer
}
