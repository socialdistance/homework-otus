package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("invalid config file", func(t *testing.T) {
		_, err := LoadConfig("/test/test.test")
		require.Error(t, err)
	})
}
