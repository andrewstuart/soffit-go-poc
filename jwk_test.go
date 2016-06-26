package main

import (
	"testing"
)

func TestGen(t *testing.T) {
	_, err := getJWT(Payload{}, "foo")

	if err != nil {
		t.Fatal(err)
	}
}
