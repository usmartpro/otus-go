package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("not exists path", func(t *testing.T) {
		tempPath := os.TempDir()
		pathName := path.Join(tempPath, "100-500-hz")
		defer os.Remove(path.Join(tempPath, "100-500-hz"))

		_, err := ReadDir(pathName)
		require.Equal(t, ErrUnknownPath, err)
	})
}
