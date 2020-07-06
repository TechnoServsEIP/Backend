package main
import "testing"

func CITEST(t *testing.T) {
	test:= 1 + 1
	if test != 3 {
		t.Errorf("1+1 = %d; want 3", test)
	}
}