package main

import (
	"testing"

	soffit "astuart.co/soffit-go"
)

func TestGen(t *testing.T) {
	_, err := getJWT(soffit.Payload{}, "foo")

	if err != nil {
		t.Fatal(err)
	}
}
