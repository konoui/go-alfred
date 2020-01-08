package alfred

import (
	"io/ioutil"
	"testing"
)

var item01 = Items{
	Item{
		Title:    "title1",
		Subtitle: "subtitle1",
	},
	Item{
		Title:    "title2",
		Subtitle: "subtitle2",
	},
}

var item02 = Items{
	Item{
		Title:    "title2",
		Subtitle: "subtitle2",
	},
	Item{
		Title:    "title1",
		Subtitle: "subtitle1",
	},
}

var item03 = Items{
	Item{
		Title:    "title3",
		Subtitle: "subtitle3",
	},
	Item{
		Title:    "title1",
		Subtitle: "subtitle1",
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
			filepath:    "./test_scriptfilter_marshal.json",
			items:       item01,
		},
		{
			description: "in the different order",
			filepath:    "./test_scriptfilter_marshal.json",
			items:       item02,
		},
		{
			description: "different values",
			filepath:    "./test_scriptfilter_marshal.json",
			items:       item03,
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
			for _, item := range tt.items {
				sf.Append(item)
			}

			gotData := sf.Marshal()
			diff := DiffScriptFilter(wantData, gotData)
			if !tt.expectedErr && diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}

		})
	}
}
