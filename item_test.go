package alfred

import (
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.want
			item := NewItem().SetTitle(input.Title).SetSubtitle(input.Subtitle).
				SetArg(input.Arg).SetAutocomplete(input.Autocomplete)

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
