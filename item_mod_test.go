package alfred

import (
	"testing"
)

func TestNewMod(t *testing.T) {
	tests := []struct {
		name string
		want *Mod
	}{
		{
			name: "new mod",
			want: &Mod{
				variables: Variables{
					"key1": "1",
					"key2": "2",
					"key3": "3",
				},
				valid:    boolP(true),
				arg:      "arg",
				subtitle: "subtitle",
				icon: &Icon{
					typ:  "typ",
					path: "path",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMod().Arg("arg").Valid(true).Subtitle("subtitle").Icon(
				NewIcon().Path("path").Type("typ"),
			).Variable("key1", "1").Variables(
				Variables{"key2": "2", "key3": "3"},
			)

			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}

			// marshal/unmarshal test
			b, err := got.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			got = NewMod()
			if err := got.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}

			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n %s", diff)
			}
		})
	}
}
