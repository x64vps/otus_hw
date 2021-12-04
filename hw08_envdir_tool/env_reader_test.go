package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("empty dir", func(t *testing.T) {
		_, err := ReadDir("")

		require.ErrorIs(t, err, ErrEmptyDirPath)
	})

	t.Run("check parse env", func(t *testing.T) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		envs, err := ReadDir(dir + "/testdata/env")
		require.NoError(t, err)

		valid := Environment{
			"UNSET": EnvValue{Value: "", NeedRemove: true},
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		}

		require.Equal(t, len(valid), len(envs))

		for k, v := range envs {
			val, ok := valid[k]
			require.True(t, ok)

			require.Equal(t, v.Value, val.Value)
			require.Equal(t, v.NeedRemove, val.NeedRemove)
		}
	})
}
