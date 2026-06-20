package svc

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
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
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			_, err := New("dummy-name", "dummy-version", tc.serviceOption)
			require.NoError(t, err)
		})
	}
}
