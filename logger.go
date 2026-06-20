package svc

import (
	"context"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func (s *SVC) initLogger(handler slog.Handler, levelVar *slog.LevelVar) {
	if s.metricsCounter != nil {
		handler = &metricsHandler{inner: handler, counter: s.metricsCounter}
	}
	s.logger = slog.New(handler)
	s.levelVar = levelVar
	s.stdLogger = slog.NewLogLogger(handler, slog.LevelError)
}

// WithSlogMetrics adds a hook to the logger and emits metrics to prometheus
// based on log level and log name.
func WithSlogMetrics() Option {
	return func(s *SVC) error {
		requestCount := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "logger_emitted_entries_total",
				Help: "Number of log messages emitted.",
			},
			[]string{"level", "logger_name"},
		)
		if err := prometheus.Register(requestCount); err != nil {
			return err
		}

		s.metricsCounter = requestCount
		return nil
	}
}

// WithLogger is an option that allows you to provide your own customized logger.
func WithLogger(logger *slog.Logger, levelVar *slog.LevelVar) Option {
	return func(s *SVC) error {
		s.initLogger(logger.Handler(), levelVar)
		return nil
	}
}

// WithDevelopmentLogger is an option that uses a JSON handler with
// configurations set meant to be used for development.
func WithDevelopmentLogger() Option {
	return func(s *SVC) error {
		levelVar := &slog.LevelVar{}
		levelVar.Set(slog.LevelDebug)

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: levelVar,
		})
		s.initLogger(handler, levelVar)
		s.logger = s.logger.With("app", s.Name, "version", s.Version)
		return nil
	}
}

// WithProductionLogger is an option that uses a JSON handler with
// configurations set meant to be used for production.
func WithProductionLogger() Option {
	return func(s *SVC) error {
		levelVar := &slog.LevelVar{}
		levelVar.Set(slog.LevelInfo)

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: levelVar,
		})
		s.initLogger(handler, levelVar)
		s.logger = s.logger.With("app", s.Name, "version", s.Version)
		return nil
	}
}

// WithConsoleLogger is an option that uses a text handler with configurations
// set meant to be used for debugging in the console.
func WithConsoleLogger(level slog.Level) Option {
	return func(s *SVC) error {
		levelVar := &slog.LevelVar{}
		levelVar.Set(level)

		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: levelVar,
		})
		s.initLogger(handler, levelVar)
		return nil
	}
}

// gcpReplaceAttr maps slog keys to GCP Stackdriver-compatible keys.
func gcpReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		a.Key = "severity"
		level := a.Value.Any().(slog.Level)
		switch {
		case level < slog.LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < slog.LevelWarn:
			a.Value = slog.StringValue("INFO")
		case level < slog.LevelError:
			a.Value = slog.StringValue("WARNING")
		default:
			a.Value = slog.StringValue("ERROR")
		}
		return a
	}
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}
	if a.Key == slog.TimeKey {
		a.Key = "timestamp"
	}
	return a
}

// WithStackdriverLogger is an option that uses a JSON handler with
// configurations set meant to be used for production and is compliant with
// the GCP/Stackdriver format.
func WithStackdriverLogger(level slog.Level) Option {
	return func(s *SVC) error {
		levelVar := &slog.LevelVar{}
		levelVar.Set(level)

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       levelVar,
			ReplaceAttr: gcpReplaceAttr,
		})
		s.initLogger(handler, levelVar)
		s.logger = s.logger.With(
			slog.Group("serviceContext", slog.String("service", s.Name)),
			slog.Group("logging.googleapis.com/labels", slog.String("version", s.Version)),
		)
		return nil
	}
}


// metricsHandler wraps a slog.Handler to count log entries by level.
type metricsHandler struct {
	inner   slog.Handler
	counter *prometheus.CounterVec
}

func (h *metricsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *metricsHandler) Handle(ctx context.Context, r slog.Record) error {
	h.counter.WithLabelValues(r.Level.String()).Inc()
	return h.inner.Handle(ctx, r)
}

func (h *metricsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &metricsHandler{h.inner.WithAttrs(attrs), h.counter}
}

func (h *metricsHandler) WithGroup(name string) slog.Handler {
	return &metricsHandler{h.inner.WithGroup(name), h.counter}
}
