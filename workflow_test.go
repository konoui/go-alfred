package alfred

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
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
			name: "Add return 2",
			fields: fields{
				s: NewScriptFilter(),
			},
			args: args{
				i: 2,
			},
			want: &Workflow{
				std: &ScriptFilter{
					rerun: 2,
				},
				warn: &ScriptFilter{
					rerun: 2,
				},
				err: &ScriptFilter{
					rerun: 2,
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
			if got := w.Rerun(tt.args.i); !reflect.DeepEqual(got, tt.want) {
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
				title: "test",
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
			filepath:    testFilePath("test_scriptfilter_items01.json"),
			items:       items01,
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
			awf.SetEmptyWarning(tt.emptyItem.title, tt.emptyItem.subtitle)
			awf.Append(tt.items...)

			got := awf.Marshal()
			if diff := DiffScriptFilter(want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestOutput(t *testing.T) {
	tests := []struct {
		description     string
		filepath        string
		scriptfilter    *ScriptFilter
		emptyItem       Item
		additionalKey   string
		additionalValue string
	}{
		{
			description:     "output standard items",
			filepath:        testFilePath("test_scriptfilter01_additional_env.json"),
			scriptfilter:    scriptfilter01,
			emptyItem:       emptyItem,
			additionalKey:   additionalKey,
			additionalValue: additionalValue,
		},
		{
			description:  "output empty warning",
			filepath:     testFilePath("test_scriptfilter_empty_warning_marshal.json"),
			scriptfilter: &ScriptFilter{},
			emptyItem:    emptyItem,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			wantData, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			awf := NewWorkflow()
			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			awf.SetOut(outBuf)
			awf.SetErr(errBuf)
			awf.SetEmptyWarning(tt.emptyItem.title, tt.emptyItem.subtitle)
			awf.Append(tt.scriptfilter.items...).
				Variables(tt.scriptfilter.variables).
				Variable(tt.additionalKey, tt.additionalValue)

			awf.Output()
			gotData := outBuf.Bytes()
			if diff := DiffScriptFilter(wantData, gotData); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}

			gotString := errBuf.String()
			if gotString != "" {
				t.Error("gotString should be empty")
			}

			awf.Output()
			gotString = errBuf.String()
			wantString := sentMessage
			if strings.Contains(wantString, gotString) {
				t.Errorf("\nwant: %s\ngot: %s\n", wantString, gotString)
			}
		})
	}
}
