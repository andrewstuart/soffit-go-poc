package main

import (
	"testing"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
)

func TestGen(t *testing.T) {
	_, err := getJWT(soffit.Request{}, "foo")

	if err != nil {
		t.Fatal(err)
	}
}
