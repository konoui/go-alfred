package alfred

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/konoui/go-alfred/logger"
)

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
				logger: logger.New(ioutil.Discard),
				dirs:   make(map[string]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewWorkflow()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: %#v, got: %#v", tt.want, got)
			}
		})
	}
}

func TestWorkflow_Rerun(t *testing.T) {
	type fields struct {
		s *ScriptFilter
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
			name: "set return 2",
			fields: fields{
				s: NewScriptFilter(),
			},
			args: args{
				i: 2,
			},
			want: &Workflow{
				std: &ScriptFilter{
					Rerun: 2,
				},
				warn: &ScriptFilter{
					Rerun: 2,
				},
				err: &ScriptFilter{
					Rerun: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				std:  tt.fields.s,
				warn: tt.fields.s,
				err:  tt.fields.s,
			}
			if got := w.SetRerun(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Workflow.Rerun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflow_Clear(t *testing.T) {
	tests := []struct {
		name string
		item *Item
		want []byte
	}{
		{
			name: "clear items",
			item: &Item{
				Title: "test",
			},
			want: NewWorkflow().Marshal(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWorkflow().Append(tt.item).Clear().Marshal()
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
			filepath:    testFilePath("test_scriptfilter_marshal.json"),
			items:       item01,
			emptyItem:   emptyItem,
		},
		{
			description: "output empty warning",
			filepath:    testFilePath("test_scriptfilter_empty_warning_marshal.json"),
			items:       Items{},
			emptyItem:   emptyItem,
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
