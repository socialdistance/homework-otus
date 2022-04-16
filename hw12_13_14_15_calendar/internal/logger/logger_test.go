package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("test debug", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		require.Nil(t, err)

		defer os.Remove(file.Name())
		defer file.Close()

		logg, err := New(config.LoggerConf{
			Level:    config.Debug,
			Filename: file.Name(),
		})
		require.Nil(t, err)

		logg.Debug("DEBUG %s", "debug")
		logg.Info("INFO %s", "info")
		logg.Warn("WARN %s", "warn")
		logg.Error("ERROR %s", "error")

		loggerRes, _ := os.ReadFile(file.Name())
		require.True(t, strings.Contains(string(loggerRes), "DEBUG debug"))
		require.True(t, strings.Contains(string(loggerRes), "INFO info"))
		require.True(t, strings.Contains(string(loggerRes), "WARN warn"))
		require.True(t, strings.Contains(string(loggerRes), "ERROR error"))
	})

	t.Run("test info", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		require.Nil(t, err)

		defer os.Remove(file.Name())
		defer file.Close()

		logg, _ := New(config.LoggerConf{
			Level:    config.Info,
			Filename: file.Name(),
		})

		logg.Debug("DEBUG %s", "debug")
		logg.Info("INFO %s", "info")
		logg.Warn("WARN %s", "warn")
		logg.Error("ERROR %s", "error")

		loggerRes, _ := os.ReadFile(file.Name())
		require.False(t, strings.Contains(string(loggerRes), "DEBUG debug"))
		require.True(t, strings.Contains(string(loggerRes), "INFO info"))
		require.True(t, strings.Contains(string(loggerRes), "WARN warn"))
		require.True(t, strings.Contains(string(loggerRes), "ERROR error"))
	})

	t.Run("test warning", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		require.Nil(t, err)

		defer os.Remove(file.Name())
		defer file.Close()

		logg, _ := New(config.LoggerConf{
			Level:    config.Warn,
			Filename: file.Name(),
		})

		logg.Debug("DEBUG %s", "debug")
		logg.Info("INFO %s", "info")
		logg.Warn("WARN %s", "warn")
		logg.Error("ERROR %s", "error")

		loggerRes, _ := os.ReadFile(file.Name())
		require.False(t, strings.Contains(string(loggerRes), "DEBUG debug"))
		require.False(t, strings.Contains(string(loggerRes), "INFO info"))
		require.True(t, strings.Contains(string(loggerRes), "WARN warn"))
		require.True(t, strings.Contains(string(loggerRes), "ERROR error"))
	})

	t.Run("test error", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		require.Nil(t, err)

		defer os.Remove(file.Name())
		defer file.Close()

		logg, _ := New(config.LoggerConf{
			Level:    config.Error,
			Filename: file.Name(),
		})

		logg.Debug("DEBUG %s", "debug")
		logg.Info("INFO %s", "info")
		logg.Warn("WARN %s", "warn")
		logg.Error("ERROR %s", "error")

		loggerRes, _ := os.ReadFile(file.Name())
		require.False(t, strings.Contains(string(loggerRes), "DEBUG debug"))
		require.False(t, strings.Contains(string(loggerRes), "INFO info"))
		require.False(t, strings.Contains(string(loggerRes), "WARN warn"))
		require.True(t, strings.Contains(string(loggerRes), "ERROR error"))
	})
}
