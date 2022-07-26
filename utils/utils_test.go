package utils

import "testing"

func TestHash(t *testing.T) {
	s := struct {
		Test string
	}{Test: "test"}
	x := Hash(s)
	t.Log(x)
}
