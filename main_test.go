package main

import "testing"

func Test(t *testing.T) {
	test := "test"
	if test != "test" {
		t.Errorf("test = %d; want test", test)
	}
}
