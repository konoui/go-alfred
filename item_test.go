package alfred

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestItemAPI(t *testing.T) {
	tests := []struct {
		name string
		want *Item
	}{
		{
			name: "chain methid is equal native",
			want: items04[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.want
			item := NewItem().
				Title(input.title).
				Subtitle(input.subtitle).
				Arg(input.arg).
				Autocomplete(input.autocomplete).
				Match(input.match).
				QuicklookURL(input.quicklookURL).
				UID(input.uid).
				Valid(input.valid).
				Icon(
					NewIcon().
						Type(input.icon.typ).
						Path(input.icon.path),
				).
				Text(
					NewText().
						Copy(input.text.copy).
						LargeType(input.text.largeType),
				)

			// TODO
			got := item.Mods(input.mods).Variables(input.variables)
			if diff := diffItemObject(tt.want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}

			// overwrite same variables with single method e.g AddVariable
			for k, v := range input.variables {
				item.Variable(k, v)
			}

			for k, v := range input.mods {
				inputMod := input.mods[k]
				mod := NewMod().Arg(inputMod.arg).Subtitle(inputMod.subtitle)
				for k, v := range inputMod.variables {
					mod.Variable(k, v)
				}
				item.Mod(k, v)
			}

			got = item
			if diff := diffItemObject(tt.want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func Test_AddVariableVariables(t *testing.T) {
	name := "Add a variable and variables"
	vals := Variables{"key1": "val1", "key2": "val2"}
	want := &Item{
		variables: Variables{
			"key0": "val0",
			"key1": "val1",
			"key2": "val2",
			"key3": "val3",
		},
	}
	t.Run(name, func(t *testing.T) {
		got := NewItem().Variable("key0", "val0").Variables(vals).Variable("key3", "val3")
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
}

func diffItemObject(want, got interface{}) string {
	out1Data, err := json.Marshal(want)
	if err != nil {
		return fmt.Sprintf("failed to marshal want due to %v", err)
	}

	out2Data, err := json.Marshal(got)
	if err != nil {
		return fmt.Sprintf("failed to marshal got due to %v", err)
	}

	return cmp.Diff(string(out1Data), string(out2Data))
}
