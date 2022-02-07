package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("negative offset", func(t *testing.T) {
		testFile, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(testFile.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", -100500, 100)
		require.Equal(t, ErrOffsetValue, err)
		testFile.Close()
	})

	t.Run("negative limit", func(t *testing.T) {
		testFile, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(testFile.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", 100, -100500)
		require.Equal(t, ErrLimitValue, err)
		testFile.Close()
	})

	t.Run("big offset", func(t *testing.T) {
		testFile, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(testFile.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", 1000000000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
		testFile.Close()
	})
}
