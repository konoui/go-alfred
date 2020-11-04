package alfred

import (
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
			name: "chaine methid is equal native",
			want: &Item{
				Title:        "title",
				Subtitle:     "subtitle",
				Arg:          "arg",
				Autocomplete: "autocomplete",
				Variables: map[string]string{
					"key": "value",
				},
				Icon: &Icon{
					Type: "image",
					Path: "./",
				},
				Mods: map[ModKey]*Mod{
					ModCtrl: {
						Subtitle: "modctrl",
						Arg:      "arg",
						Variables: map[string]string{
							"key": "value",
						},
					},
				},
				Text: &Text{
					Copy:      "copy",
					LargeType: "largetype",
				},
				Match:        "match-value",
				QuicklookURL: "quick-url",
				UID:          "uid",
				Valid:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.want
			item := NewItem().SetTitle(input.Title).SetSubtitle(input.Subtitle).
				SetArg(input.Arg).SetAutocomplete(input.Autocomplete).
				SetMatch(input.Match).SetQuicklookURL(input.QuicklookURL).
				SetUID(input.UID).SetValid(input.Valid).
				SetMods(input.Mods).SetMod(ModCtrl, input.Mods[ModCtrl])

			for k, v := range input.Variables {
				item.SetVariable(k, v)
			}

			for k, v := range input.Mods {
				inputMod := input.Mods[k]
				mod := NewMod().SetArg(inputMod.Arg).SetSubtitle(inputMod.Subtitle)
				for k, v := range inputMod.Variables {
					mod.SetVariable(k, v)
				}
				item.SetMod(k, v)
			}

			icon := NewIcon().SetType(input.Icon.Type).SetPath(input.Icon.Path)
			text := NewText().SetCopy(input.Text.Copy).SetLargeType(input.Text.LargeType)
			item.SetIcon(icon).SetText(text)

			got := item
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func Test_SetVariableVariables(t *testing.T) {
	name := "set a variable and variables"
	vals := Variables{"key1": "val1", "key2": "val2"}
	want := &Item{
		Variables: Variables{
			"key0": "val0",
			"key1": "val1",
			"key2": "val2",
			"key3": "val3",
		},
	}
	t.Run(name, func(t *testing.T) {
		got := NewItem().SetVariable("key0", "val0").SetVariables(vals).SetVariable("key3", "val3")
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
}
