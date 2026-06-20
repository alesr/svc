package main

import (
	"log/slog"

	"github.com/alesr/svc"
)

var _ svc.Worker = (*dummyWorker)(nil)

type dummyWorker struct{}

func (d *dummyWorker) Init(*slog.Logger) error { return nil }
func (d *dummyWorker) Terminate() error        { return nil }
func (d *dummyWorker) Run() error              { select {} }

func main() {
	s, err := svc.New("minimal-service", "1.0.0")
	svc.MustInit(s, err)

	w := &dummyWorker{}
	s.AddWorker("dummy-worker", w)

	s.Run()
}
