package main

import (
	"testing"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func TestPossessiveComma(t *testing.T) {
	t.Logf("Testing the possesive comma func")

	tests := map[string]bool{
		"can't": true,       // true
		"he's": true,        // true
		"rock'n'roll": true, // true (multiple matches)
		"won't": true,       // true
		"it's": true,        // true
		"'start": false,      // false
		"end'": false,        // false
		"no quote": false,    // false
	}

	for test, result := range tests{
		if notPossessiveComma(test) == result{
			aphrodite.PrintColour("Green", test + "\n")
		} else {
			aphrodite.PrintColour("Red", test + "\n")
		}
	}
}
