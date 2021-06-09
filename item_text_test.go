package alfred

import (
	"testing"
)

func TestNewText(t *testing.T) {
	tests := []struct {
		name string
		want *Text
	}{
		{
			name: "new text",
			want: &Text{
				copy:      "copy",
				largeType: "large",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewText().Copy("copy").LargeType("large")
			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}

			// marshal/unmarshal test
			b, err := got.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			got = NewText()
			if err := got.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}

			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}
		})
	}
}
