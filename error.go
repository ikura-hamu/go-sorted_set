package gosortedset

import (
	"cmp"
	"errors"
)

var (
	ErrIndexOutOfRange = errors.New("index out of range")
)

func Must[T cmp.Ordered](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
