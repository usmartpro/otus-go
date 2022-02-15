package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("ls", func(t *testing.T) {
		environment := make(map[string]EnvValue)
		code := RunCmd([]string{"ls", "-la"}, environment)
		require.Equal(t, 0, code)
	})
}
