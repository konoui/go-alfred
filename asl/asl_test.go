package asl

import (
	"bytes"
	"io"
	"testing"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "write simple data",
			input: "test-data",
		},
		{
			name:  "write simple return string",
			input: "test-data\ntest-data\n",
		},
		{
			name:  "a * 10000",
			input: repeat(10000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asl := New()
			out := &bytes.Buffer{}
			mw := io.MultiWriter(asl, out)
			n, err := mw.Write([]byte(tt.input))
			if err != nil {
				t.Error(err)
			}

			if got := len(tt.input); got != n {
				t.Errorf("number of byte want: %v, got %v", n, got)
			}

			got := out.String()
			if got != tt.input {
				t.Errorf("unexpected string want: %v, got %v", tt.input, got)
			}
		})
	}
}

func repeat(n int) string {
	s := "a"
	for i := 0; i < n; i++ {
		s += "a"
	}
	return s
}
