package alfred

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

var scriptfilter01 = &ScriptFilter{
	items: items01,
	variables: Variables{
		"key1": "value1",
		"key2": "value2",
	},
	rerun: 2,
}
var additionalKey = "key3"
var additionalValue = "value3"

func TestNewScriptFilter(t *testing.T) {
	tests := []struct {
		description string
		want        *ScriptFilter
	}{
		{
			description: "create new scriptfilter",
			want:        &ScriptFilter{},
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
		description  string
		filepath     string
		scriptfilter *ScriptFilter
		key          string
		value        string
	}{
		{
			description:  "create new scriptfilter",
			filepath:     testFilePath("test_scriptfilter01_additional_env.json"),
			scriptfilter: scriptfilter01,
			key:          additionalKey,
			value:        additionalValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			want, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			sf := NewScriptFilter()
			sf.Append(tt.scriptfilter.items...)
			sf.Variables(tt.scriptfilter.variables)
			sf.Variable(tt.key, tt.value)

			got := sf.Marshal()
			if diff := DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
	}{
		{
			description: "unmarshal test",
			filepath:    testFilePath("test_scriptfilter_items04.json"),
			items:       items04,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			input, err := ioutil.ReadFile(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			wantSf := NewScriptFilter()
			wantSf.Append(tt.items...)

			gotSf := ScriptFilter{}
			err = json.Unmarshal(input, &gotSf)
			if err != nil {
				t.Fatal(err)
			}

			got := gotSf.Marshal()
			want := wantSf.Marshal()
			if diff := DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
