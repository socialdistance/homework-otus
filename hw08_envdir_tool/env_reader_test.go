package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read env dir", func(t *testing.T) {
		expected := Environment{
			"BAR":   EnvValue{"bar", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"EMPTY": EnvValue{"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("testdata/env")
		require.ErrorIs(t, nil, err)
		for k, v := range expected {
			require.Equal(t, v, env[k], "%v wrong value", k)
		}
	})

	t.Run("first line", func(t *testing.T) {
		res, err := checkFirstLine("./testdata/env/FOO")
		require.Nil(t, err)
		require.Equal(t, "   foo\nwith new line", clear(res))
	})

	t.Run("exist file", func(t *testing.T) {
		_, err := ReadDir("./testdata/env/test")
		require.Error(t, err)
	})
}
