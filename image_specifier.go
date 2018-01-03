package main

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/thraxil/resize"
)

// combination of fields that uniquely specify an image
type imageSpecifier struct {
	Hash      *hash
	Size      *resize.SizeSpec
	Extension string // with leading '.'
}

func (i imageSpecifier) String() string {
	return path.Join(i.Hash.String(), i.Size.String(), "image"+i.Extension)
}

func resizedPath(imagePath, size string) string {
	d := filepath.Dir(imagePath)
	extension := filepath.Ext(imagePath)
	return path.Join(d, size+extension)
}

func newImageSpecifier(s string) *imageSpecifier {
	parts := strings.Split(s, "/")
	ahash, _ := hashFromString(parts[0], "")
	size := parts[1]
	rs := resize.MakeSizeSpec(size)
	filename := parts[2]
	fparts := strings.Split(filename, ".")
	extension := "." + fparts[1]
	return &imageSpecifier{Hash: ahash, Size: rs, Extension: extension}
}

func (i imageSpecifier) sizedPath() string {
	return resizedPath(i.fullSizePath(), i.Size.String())
}

func (i imageSpecifier) baseDir() string {
	return i.Hash.String()
}

func (i imageSpecifier) fullSizePath() string {
	return path.Join(i.baseDir(), "full"+i.Extension)
}

func (i imageSpecifier) fullVersion() imageSpecifier {
	i.Size = resize.MakeSizeSpec("full")
	return i
}
