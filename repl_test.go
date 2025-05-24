package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "Sneasel",
			expected: []string{"sneasel"},
		},
		{
			input:    "          		",
			expected: []string{},
		},
		{
			input:    "Chikorita	 cyndaQUIL    totodile",
			expected: []string{"chikorita", "cyndaquil", "totodile"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check the length of the actual slice against the expected slice. If they don't match, use t.Errorf to print an error message and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("TEST FAILED: length mismatch")
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice. If they don't match, use t.Errorf to print an error message and fail the test
			if word != expectedWord {
				t.Errorf("TEST FAILED: letter mismatch")
				break
			}
		}
	}
}
