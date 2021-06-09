package alfred

import (
	"testing"
)

func TestNewIcon(t *testing.T) {
	tests := []struct {
		name string
		want *Icon
	}{
		{
			name: "new icon",
			want: &Icon{
				typ:  "typ",
				path: "path",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIcon().Path("path").Type("typ")

			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}

			// marshal/unmarshal test
			b, err := got.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			got = NewIcon()
			if err := got.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}

			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}
		})
	}
}
