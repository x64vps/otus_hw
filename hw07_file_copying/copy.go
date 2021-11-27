package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitLessThanZero     = errors.New("limit less than zero")
	ErrOffsetLessThanZero    = errors.New("offset less than zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validate(fromPath, offset, limit); err != nil {
		return err
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if _, err = fromFile.Seek(offset, 0); err != nil {
		return err
	}

	if limit == 0 {
		fi, _ := fromFile.Stat()
		limit = fi.Size()
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barReader, limit)

	bar.Finish()

	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

func validate(fromPath string, offset, limit int64) error {
	if limit < 0 {
		return ErrLimitLessThanZero
	}

	if offset < 0 {
		return ErrOffsetLessThanZero
	}

	fi, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}
