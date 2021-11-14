package alfred

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

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
	}{
		{
			name: "Add return 2",
			fields: fields{
				s: NewScriptFilter(),
			},
			args: args{
				i: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow().Rerun(tt.args.i)
			if !(w.err.rerun == tt.args.i && w.std.rerun == tt.args.i && w.warn.rerun == tt.args.i) {
				t.Errorf("Workflow.Rerun() = %v, want %v", w, tt.args.i)
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
			want: NewWorkflow().Bytes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWorkflow().Append(tt.item).Clear().Bytes()
			if diff := DiffOutput(tt.want, got); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}

func TestWorfkfloByte(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
		emptyItem   *Item
		systemItem  *Item
		opts        []Option
	}{
		{
			description: "output standard items",
			filepath:    testFilePath("test_scriptfilter_items01.json"),
			items:       items01,
			emptyItem:   &emptyItem,
		},
		{
			description: "output empty warning",
			filepath:    testFilePath("test_scriptfilter_empty_warning_marshal.json"),
			items:       Items{},
			emptyItem:   &emptyItem,
		},
		{
			description: "limit standard items",
			filepath:    testFilePath("test_limit_output_item01.json"),
			items:       items01,
			emptyItem:   &emptyItem,
			opts: []Option{
				WithMaxResults(1),
			},
		},
		{
			description: "set system item and  limit standard items",
			filepath:    testFilePath("test_system_and_limit_output_item01.json"),
			items:       items01,
			emptyItem:   &emptyItem,
			systemItem:  &systemItem,
			opts: []Option{
				WithMaxResults(1),
			},
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

			awf := NewWorkflow(tt.opts...)
			awf.SetEmptyWarning(tt.emptyItem.title, tt.emptyItem.subtitle)
			awf.SetSystemInfo(tt.systemItem)
			awf.Append(tt.items...)

			got := awf.Bytes()
			if diff := DiffOutput(want, got); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}

			if diff := Diff(tt.items, awf.std.items); diff != "" {
				t.Errorf("limit does not work -want +got\n%+v", diff)
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

			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			awf := NewWorkflow(
				WithOutWriter(outBuf),
				WithLogWriter(errBuf),
			)
			awf.SetEmptyWarning(tt.emptyItem.title, tt.emptyItem.subtitle)
			awf.Append(tt.scriptfilter.items...).
				Variables(tt.scriptfilter.variables).
				Variable(tt.additionalKey, tt.additionalValue)

			awf.Output()
			gotData := outBuf.Bytes()
			if diff := DiffOutput(wantData, gotData); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
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

func TestWorkflow_Fatal(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "fatal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow()
			osExit = func(code int) {}
			w.Fatal("test", "test1")
		})
	}
}
