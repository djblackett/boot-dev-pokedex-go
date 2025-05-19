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
			input:    "  foo bar baz  ",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input:    "  a B c D e  ",
			expected: []string{"a", "b", "c", "d", "e"},
			// add more cases here
		},
		{
			input:    "   Moose     Meat   ",
			expected: []string{"moose", "meat"},
		},
	}
	for _, c := range cases {
		actual := CleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected %s, but got %s", expectedWord, word)
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
