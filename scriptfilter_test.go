package alfred

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNewScriptFilter(t *testing.T) {
	tests := []struct {
		description string
		want        ScriptFilter
	}{
		{
			description: "create new workflow",
			want:        ScriptFilter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewScriptFilter()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}

		})
	}
}

func TestScriptFilterMarshal(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
	}{
		{
			description: "create new scriptfilter",
			filepath:    testFilePath("test_scriptfilter_marshal.json"),
			items:       item01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			want, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			wf := NewScriptFilter()
			for _, item := range tt.items {
				wf.Append(item)
			}

			got := wf.Marshal()
			if diff := DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
