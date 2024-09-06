package main

import (
	"reflect"
	"testing"
)

func TestCharFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:  "Basic case",
			input: "aaabbc",
			expected: map[string]int{
				"'a'": 3,
				"'b'": 2,
				"'c'": 1,
			},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:  "With spaces and special characters",
			input: "hello world! 123",
			expected: map[string]int{
				"'h'": 1, "'e'": 1, "'l'": 3, "'o'": 2, "' '": 2, "'w'": 1,
				"'r'": 1, "'d'": 1, "'!'": 1, "'1'": 1, "'2'": 1, "'3'": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charMap := make(map[string]int)
			for _, v := range tt.input {
				counter(v, charMap)
			}
			if !reflect.DeepEqual(charMap, tt.expected) {
				t.Errorf("countFrequency() = %v, want %v", charMap, tt.expected)
			}
		})
	}
}
