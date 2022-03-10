package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	t.Run("exit code", func(t *testing.T) {
		command := dir + "/testdata/test_code_100.sh"

		code, err := RunCmd([]string{command}, nil)

		require.NoError(t, err)
		require.Equal(t, 100, code)
	})

	t.Run("empty commands", func(t *testing.T) {
		_, err := RunCmd([]string{}, nil)

		require.ErrorIs(t, ErrEmptyCommand, err)
	})
}
