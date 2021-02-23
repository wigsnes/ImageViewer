package main

import "testing"

func TestRemoveLastPath(t *testing.T) {
	result := removeLastPath("this/is/a/path/")
	if result != "this/is/a/" {
		t.Errorf("Failed - Expected: %v, got: %v", "this/is/a/", result)
	}

	result2 := removeLastPath("this/")
	if result2 != "/" {
		t.Errorf("Failed - Expected: %v, got: %v", "/", result2)
	}
}
