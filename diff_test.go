package alfred

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

var emptyItem = Item{
	title:    "emptyTitle1",
	subtitle: "emptySubtitle1",
}

var items01 = Items{
	&Item{
		title:    "title1",
		subtitle: "subtitle1",
	},
	&Item{
		title:    "title2",
		subtitle: "subtitle2",
	},
}

var items02 = Items{
	&Item{
		title:    "title2",
		subtitle: "subtitle2",
	},
	&Item{
		title:    "title1",
		subtitle: "subtitle1",
	},
}

var items03 = Items{
	&Item{
		title:    "title3",
		subtitle: "subtitle3",
	},
	&Item{
		title:    "title1",
		subtitle: "subtitle1",
	},
}

var items04 = Items{
	&Item{
		title:        "title",
		subtitle:     "subtitle",
		arg:          "arg",
		autocomplete: "autocomplete",
		variables: map[string]string{
			"key": "value",
		},
		icon: &Icon{
			typ:  "image",
			path: "./",
		},
		mods: map[ModKey]*Mod{
			ModCtrl: {
				subtitle: "modctrl",
				arg:      "arg",
				variables: map[string]string{
					"key": "value",
				},
				valid: boolP(false),
			},
		},
		text: &Text{
			copy:      "copy",
			largeType: "largetype",
		},
		match:        "match-value",
		quicklookURL: "quick-url",
		uid:          "uid",
		valid:        boolP(false),
	},
}

func TestDiffScriptFilter(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
		expectedErr bool
	}{
		{
			description: "in the same order",
			filepath:    testFilePath("test_scriptfilter_items01.json"),
			items:       items01,
		},
		{
			description: "in the different order",
			filepath:    testFilePath("test_scriptfilter_items01.json"),
			items:       items02,
		},
		{
			description: "different values",
			filepath:    testFilePath("test_scriptfilter_items01.json"),
			items:       items03,
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			wantData, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			sf := NewScriptFilter()
			sf.Append(tt.items...)

			gotData := sf.Marshal()
			diff := DiffOutput(wantData, gotData)
			if !tt.expectedErr && diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}

func testFilePath(f string) string {
	return filepath.Join("testdata", f)
}

func boolP(v bool) *bool {
	return &v
}
