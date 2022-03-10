package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

var ErrEmptyDirPath = errors.New("dir is empty")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if dir == "" {
		return nil, ErrEmptyDirPath
	}

	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, d := range dirs {
		v, err := getEnvValue(dir, d)
		if err != nil {
			return nil, err
		}

		env[d.Name()] = *v
	}

	return env, nil
}

func getEnvValue(path string, dir os.DirEntry) (*EnvValue, error) {
	f, err := os.Open(path + "/" + dir.Name())
	if err != nil {
		return nil, err
	}

	fi, _ := f.Stat()
	if fi.Size() == 0 {
		return &EnvValue{Value: "", NeedRemove: true}, nil
	}

	reader := bufio.NewReader(f)
	l, _, err := reader.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	l = bytes.ReplaceAll(l, []byte{'\000'}, []byte{'\n'})

	s := strings.TrimRight(string(l), " \t")

	return &EnvValue{Value: s, NeedRemove: s == ""}, nil
}
