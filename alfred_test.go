package alfred

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testItems = Items{
	&Item{
		Title:    "title1",
		Subtitle: "subtitle1",
	},
	&Item{
		Title:    "title2",
		Subtitle: "subtitle2",
	},
}

var testEmptyItem = Item{
	Title:    "emptyTitle1",
	Subtitle: "emptySubtitle1",
}

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
			filepath:    "./test_scriptfilter_marshal.json",
			items:       testItems,
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

func TestNewWorkflow(t *testing.T) {
	tests := []struct {
		description string
		want        *Workflow
	}{
		{
			description: "create new workflow",
			want: &Workflow{
				std:  NewScriptFilter(),
				warn: NewScriptFilter(),
				err:  NewScriptFilter(),
				streams: streams{
					out: os.Stdout,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewWorkflow()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}

		})
	}
}

func TestWorfkflowMarshal(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
		emptyItem   Item
	}{
		{
			description: "output standard items",
			filepath:    "./test_scriptfilter_marshal.json",
			items:       testItems,
			emptyItem:   testEmptyItem,
		},
		{
			description: "output empty warning",
			filepath:    "./test_scriptfilter_empty_warning_marshal.json",
			items:       Items{},
			emptyItem:   testEmptyItem,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			want, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			awf := NewWorkflow()
			awf.EmptyWarning(tt.emptyItem.Title, tt.emptyItem.Subtitle)
			for _, item := range tt.items {
				awf.Append(item)
			}

			got := awf.Marshal()
			if diff := DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
