package main

import "errors"

type hash struct {
	Algorithm string
	Value     []byte
}

func hashFromString(str, algorithm string) (*hash, error) {
	if algorithm == "" {
		algorithm = "sha1"
	}
	if len(str) != 40 {
		return nil, errors.New("invalid hash")
	}
	return &hash{algorithm, []byte(str)}, nil
}

func (h hash) String() string {
	return string(h.Value)
}

func (h hash) Valid() bool {
	return h.Algorithm == "sha1" && len(h.String()) == 40
}
