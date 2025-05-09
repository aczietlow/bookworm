package main

import "testing"

func textCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: " testing white space ",
			expected: []string{"hello", "world"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleaned input contained length of %v, was expected %v", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("received %s, was expecting %s", word, expectedWord)
			}
		}

	}
}
