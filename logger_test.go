package svc

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		serviceOption Option
	}{
		{
			name:          "console logger",
			serviceOption: WithConsoleLogger(slog.LevelInfo),
		},
		{
			name:          "development logger",
			serviceOption: WithDevelopmentLogger(),
		},
		{
			name:          "production logger",
			serviceOption: WithProductionLogger(),
		},
		{
			name:          "stackdriver logger",
			serviceOption: WithStackdriverLogger(slog.LevelWarn),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := New("dummy-name", "dummy-version", tt.serviceOption)
			require.NoError(t, err)
		})
	}
}
