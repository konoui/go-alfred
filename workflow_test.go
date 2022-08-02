package alfred

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func testWorkflow(opts ...Option) *Workflow {
	d := append([]Option{WithLogWriter(io.Discard), WithOutWriter(io.Discard)}, opts...)
	return NewWorkflow(d...)
}

func TestWorkflow_Rerun(t *testing.T) {
	type args struct{ i Rerun }
	tests := []struct {
		name         string
		scriptFilter *ScriptFilter
		args         args
	}{
		{
			name:         "Add return 2",
			scriptFilter: NewScriptFilter(),
			args:         args{i: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow().Rerun(tt.args.i)
			if !(w.rerun == tt.args.i) {
				t.Errorf("Workflow.Rerun() = %v, want %v", w, tt.args.i)
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
			opts:        []Option{WithMaxResults(1)},
		},
		{
			description: "output empty warning and system info",
			filepath:    testFilePath("test_scriptfilter_empty_warning_system_info_marshal.json"),
			items:       Items{},
			emptyItem:   &emptyItem,
			systemItem:  &systemItem,
		},
		{
			description: "set system item and  limit standard items",
			filepath:    testFilePath("test_system_and_limit_output_item01.json"),
			items:       items01,
			emptyItem:   &emptyItem,
			systemItem:  &systemItem,
			opts:        []Option{WithMaxResults(1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			want, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			awf := testWorkflow(tt.opts...)
			awf.SetEmptyWarning(tt.emptyItem.title, tt.emptyItem.subtitle)
			awf.SetSystemInfo(tt.systemItem)
			awf.Append(tt.items...)

			got := awf.Bytes()
			if diff := DiffOutput(want, got); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}

			if diff := Diff(tt.items, awf.items); diff != "" {
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

			wantData, err := io.ReadAll(f)
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

func TestWorkflow_Clear(t *testing.T) {
	t.Run("clear item", func(t *testing.T) {
		want := NewWorkflow().Bytes()
		got := NewWorkflow().Append(NewItem().Title("title")).Clear().Bytes()
		if diff := DiffOutput(want, got); diff != "" {
			t.Errorf("-want +got\n%+v", diff)
		}
	})
}

func TestWorkflow_Fatal(t *testing.T) {
	t.Run("fatal", func(t *testing.T) {
		osExit = func(code int) {}
		w := testWorkflow()
		w.Fatal("title", "subtitle")
		got := w.Bytes()

		item := NewItem().
			Title("title").
			Subtitle("subtitle").
			Valid(false).
			Icon(w.Asseter().IconCaution())
		want := NewWorkflow().Append(item).Bytes()
		if diff := DiffOutput(want, got); diff != "" {
			t.Errorf("-want +got\n%+v", diff)
		}
	})
}
