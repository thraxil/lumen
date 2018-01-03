package main

import (
	"fmt"
	"io"
)

type backend interface {
	fmt.Stringer
	Write(string, io.ReadCloser) error
	Read(string) ([]byte, error)
	Exists(string) bool
	Delete(string) error
	FreeSpace() uint64
}
