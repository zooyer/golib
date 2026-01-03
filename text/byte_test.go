package text

import "testing"

func TestIsPlainText(t *testing.T) {
	var tests = []struct {
		Input  []byte
		Result bool
	}{
		{Input: []byte(""), Result: true},
		{Input: []byte("\r\n\t\f"), Result: true},
		{Input: []byte("Hello World"), Result: true},
		{Input: []byte("\x00"), Result: false},
		{Input: []byte("\x82"), Result: false},
		{Input: []byte("Hello World\x00"), Result: false},
		{Input: []byte("Hello \x01\x56"), Result: false},
	}

	for i, test := range tests {
		if result := IsPlainText(test.Input); result != test.Result {
			t.Fatal(i, string(test.Input))
		}
	}
}
