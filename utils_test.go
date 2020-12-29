package alfred

import (
	"testing"

	"golang.org/x/text/unicode/norm"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "convert",
			query: norm.NFD.String("ブックマーク"),
			want:  "ブックマーク",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == tt.query {
				t.Fatal("invalid test")
			}
			if got := Normalize(tt.query); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
