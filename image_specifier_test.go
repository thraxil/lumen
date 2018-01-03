package main

import (
	_ "fmt"
	"testing"

	"github.com/thraxil/resize"
)

func Test_Create(t *testing.T) {
	s := "112e42f26fce70d268438ac8137d81607499ee10/200s/1250.jpg"
	i := newImageSpecifier(s)
	if i.Extension != ".jpg" {
		t.Error("wrong extension")
	}
}

func Test_String(t *testing.T) {
	s := "112e42f26fce70d268438ac8137d81607499ee10/200s/1250.jpg"
	i := newImageSpecifier(s)
	if i.String() != "112e42f26fce70d268438ac8137d81607499ee10/200s/image.jpg" {
		t.Error("incorrect stringification")
	}
}

func Test_FullSizePath(t *testing.T) {
	s := "112e42f26fce70d268438ac8137d81607499ee10/200s/1250.jpg"
	i := newImageSpecifier(s)
	r := i.fullSizePath()
	if r != "112e42f26fce70d268438ac8137d81607499ee10/full.jpg" {
		t.Errorf("wrong fullSizePath: %s", r)
	}
}

func Test_SizedPath(t *testing.T) {
	s := "112e42f26fce70d268438ac8137d81607499ee10/200s/1250.jpg"
	i := newImageSpecifier(s)
	r := i.sizedPath()
	if r != "112e42f26fce70d268438ac8137d81607499ee10/200s.jpg" {
		t.Errorf("wrong sizedPath: %s", r)
	}

	ahash, _ := hashFromString("112e42f26fce70d268438ac8137d81607499ee10", "")
	i = &imageSpecifier{
		ahash,
		resize.MakeSizeSpec("full"),
		".jpg",
	}

	r = i.sizedPath()
	if r != "112e42f26fce70d268438ac8137d81607499ee10/full.jpg" {
		t.Errorf("wrong sizedPath: %s", r)
	}

}
