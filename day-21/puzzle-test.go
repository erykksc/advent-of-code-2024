package main

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		code     string
		expected []rune
	}{
		{"029A", []rune("")},
		{"123", []rune("")},
		{"456", []rune("")},
		{"789", []rune("")},
	}

	for _, test := range tests {
		result := solve(test.code)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("solve(%q) = %q; expected %q", test.code, result, test.expected)
		}
	}
}
