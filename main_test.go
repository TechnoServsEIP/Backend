package main

import "testing"

func TestDemo(t *testing.T) {
	a := "test"
	if a != "test" {
		t.Errorf("a is different from test")
	}
}
