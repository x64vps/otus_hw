package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	f, err := os.CreateTemp("", "hw7_test_copy")
	if err != nil {
		panic(err)
	}

	t.Run("offset less than zero", func(t *testing.T) {
		err := Copy("testdata/input.txt", f.Name(), -1, 0)
		require.ErrorIs(t, err, ErrOffsetLessThanZero)
	})

	t.Run("limit less than zero", func(t *testing.T) {
		err := Copy("testdata/input.txt", f.Name(), 0, -1)
		require.ErrorIs(t, err, ErrLimitLessThanZero)
	})

	t.Run("offset greater than file", func(t *testing.T) {
		err := Copy("testdata/input.txt", f.Name(), 999999, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("not exists file", func(t *testing.T) {
		err := Copy("lalala.txt", f.Name(), 999999, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}

func TestCopyWithTestdata(t *testing.T) {
	tests := []struct {
		filePath string
		offset   int64
		limit    int64
	}{
		{filePath: "testdata/out_offset0_limit0.txt", offset: 0, limit: 0},
		{filePath: "testdata/out_offset0_limit10.txt", offset: 0, limit: 10},
		{filePath: "testdata/out_offset0_limit1000.txt", offset: 0, limit: 1000},
		{filePath: "testdata/out_offset0_limit10000.txt", offset: 0, limit: 10000},
		{filePath: "testdata/out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{filePath: "testdata/out_offset6000_limit1000.txt", offset: 6000, limit: 1000},
	}

	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode

	temp, err := os.CreateTemp("", "hw7_test_copy")
	if err != nil {
		panic(err)
	}
	defer temp.Close()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.filePath, func(t *testing.T) {
			err := Copy("testdata/input.txt", temp.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			equal, err := cmp.CompareFile(tc.filePath, temp.Name())
			require.NoError(t, err)
			require.True(t, equal)
		})
	}
}
