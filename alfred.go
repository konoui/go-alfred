package alfred

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const fatalErrorJSON = `{"items": [{"title": "Fatal Error","subtitle": "%s",}]}`

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() ScriptFilter {
	return ScriptFilter{}
}

// Append a new Item to Items
func (s *ScriptFilter) Append(item *Item) {
	s.Items = append(s.Items, item)
}

// Marshal ScriptFilter as Json
func (s *ScriptFilter) Marshal() []byte {
	res, err := json.Marshal(s)
	if err != nil {
		return []byte(fmt.Sprintf(fatalErrorJSON, err.Error()))
	}

	return res
}

// SetOut redirect stdout and stderr to s
func (w *Workflow) SetOut(out io.Writer) {
	w.streams.out = out
}

// NewWorkflow has simple ScriptFilter api
func NewWorkflow() *Workflow {
	wf := &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
		},
	}

	return wf
}

// Append a new Item to standard ScriptFilter
func (w *Workflow) Append(item *Item) *Workflow {
	w.std.Append(item)
	return w
}

// Rerun set rerun variable
func (w *Workflow) Rerun(i Rerun) *Workflow {
	w.std.Rerun = i
	w.warn.Rerun = i
	return w
}

// Variables set variables
func (w *Workflow) Variables(v Variables) *Workflow {
	w.std.Variables = v
	return w
}

// EmptyWarning create a new Item to Marshal　when there are no standard items
func (w *Workflow) EmptyWarning(title, subtitle string) {
	w.warn = NewScriptFilter()
	w.warn.Append(
		&Item{
			Title:    title,
			Subtitle: subtitle,
		})
}

// error append a new Item to error ScriptFilter
func (w *Workflow) error(title, subtitle string) {
	w.err = NewScriptFilter()
	w.err.Append(
		&Item{
			Title:    title,
			Subtitle: subtitle,
		})
}

// Marshal WorkFlow results
func (w *Workflow) Marshal() []byte {
	if len(w.std.Items) == 0 {
		return w.warn.Marshal()
	}

	return w.std.Marshal()
}

// Fatal output error to io stream
func (w *Workflow) Fatal(title, subtitle string) {
	w.error(title, subtitle)
	res := w.err.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
}

// Output to io stream
func (w *Workflow) Output() {
	res := w.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
}
