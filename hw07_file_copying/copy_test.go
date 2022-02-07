package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	from := "./testdata/input.txt"
	to := "./testdata/out.txt"

	t.Run("copy without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 0)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat(from)
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile(from)
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copy 10 bytes without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 10)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat("./testdata/out_offset0_limit10.txt")
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile("./testdata/out_offset0_limit10.txt")
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copy 1000 bytes without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat("./testdata/out_offset0_limit1000.txt")
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile("./testdata/out_offset0_limit1000.txt")
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copy more bytes then file size without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 10000)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat("./testdata/out_offset0_limit10000.txt")
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile("./testdata/out_offset0_limit10000.txt")
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copy 1000 bytes with 100 bytes offset", func(t *testing.T) {
		err := Copy(from, to, 100, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat("./testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile("./testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copy 1000 bytes with 6000 bytes offset", func(t *testing.T) {
		err := Copy(from, to, 6000, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		fileReadSize, err := os.Stat("./testdata/out_offset6000_limit1000.txt")
		require.NoError(t, err)
		fileWriteSize, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, fileReadSize.Size(), fileWriteSize.Size())

		contentFile, err := ioutil.ReadFile("./testdata/out_offset6000_limit1000.txt")
		require.NoError(t, err)
		copyFile, err := ioutil.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, contentFile, copyFile)
	})

	t.Run("copying with large offset", func(t *testing.T) {
		err := Copy(from, to, 1000000, 0)
		defer os.Remove(to)
		require.True(t, err != nil)
	})
}
