package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/alesr/svc"
)

var _ svc.Worker = (*dummyWorker)(nil)

type dummyWorker struct {
	state int
}

func (d *dummyWorker) Init(*slog.Logger) error { return nil }
func (d *dummyWorker) Terminate() error        { return nil }
func (d *dummyWorker) Run() error {
	time.Sleep(1 * time.Second)
	d.state = 1
	select {}
}

func (d *dummyWorker) Alive() error {
	if d.state == 1 {
		return fmt.Errorf("service not well, please restart me")
	}
	return nil
}

func main() {
	s, err := svc.New("minimal-service", "1.0.0", svc.WithHTTPServer("9090"), svc.WithHealthz())
	svc.MustInit(s, err)

	w := &dummyWorker{
		state: 0,
	}
	s.AddWorker("dummy-worker", w)

	s.Run()
}
