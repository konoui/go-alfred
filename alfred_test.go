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
					err: ioutil.Discard,
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

func TestWorkflow_Rerun(t *testing.T) {
	type fields struct {
		std ScriptFilter
	}
	type args struct {
		i Rerun
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Workflow
	}{
		{
			name: "set return 1",
			fields: fields{
				std: NewScriptFilter(),
			},
			args: args{
				i: 1,
			},
			want: &Workflow{
				std: ScriptFilter{
					Rerun: 1,
				},
				warn: ScriptFilter{
					Rerun: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				std: tt.fields.std,
			}
			if got := w.Rerun(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Workflow.Rerun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflow_Delete(t *testing.T) {
	tests := []struct {
		name string
		item *Item
		want []byte
	}{
		{
			name: "delete item",
			item: &Item{
				Title: "loading",
				UID:   "loading",
				Valid: true,
			},
			want: NewWorkflow().Append(&Item{
				Title: "loading",
				UID:   "loading",
				Valid: false,
			}).Marshal(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWorkflow().Append(tt.item).Delete(tt.item.UID).Marshal()
			if diff := DiffScriptFilter(tt.want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
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
