package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type lastRune struct {
	r      rune
	exists bool
	escape bool
}

func (l *lastRune) Get() rune {
	l.exists = false
	return l.r
}

func (l *lastRune) String() string {
	return string(l.Get())
}

func (l *lastRune) Set(r rune) {
	l.r = r

	l.exists = true
	l.escape = IsSlash(r)
}

func (l *lastRune) SetEscaped(r rune) error {
	if !unicode.IsDigit(r) && !IsSlash(r) {
		return ErrInvalidString
	}

	l.Set(r)
	l.escape = false

	return nil
}

func Unpack(str string) (string, error) {
	var last lastRune
	var result strings.Builder

	for _, r := range str {
		if last.escape {
			if err := last.SetEscaped(r); err != nil {
				return "", err
			}

			continue
		}

		if unicode.IsDigit(r) {
			if !last.exists {
				return "", ErrInvalidString
			}

			d, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}

			result.WriteString(strings.Repeat(last.String(), d))
			continue
		}

		if last.exists {
			result.WriteRune(last.Get())
		}

		last.Set(r)
	}

	if last.exists {
		result.WriteRune(last.Get())
	}

	return result.String(), nil
}

func IsSlash(r rune) bool {
	return string(r) == `\`
}
