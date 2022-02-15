package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test bash script", func(t *testing.T) {
		script := RunCmd([]string{"bash", "test.sh"}, Environment{})
		require.Equal(t, 0, script)
	})
}
