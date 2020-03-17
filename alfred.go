package alfred

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// see https://www.alfredapp.com/help/workflows/inputs/script-filter/json/

// Icon displayed in the result row
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// ModKey is a mod key pressed by the user to run an alternate
type ModKey string

// Valid attribute to mark if the result is valid based on the modifier selection and set a different arg to be passed out if actioned with the modifier.
const (
	ModCmd   ModKey = "cmd"   // Alternate action for ⌘↩
	ModAlt   ModKey = "alt"   // Alternate action for ⌥↩
	ModOpt   ModKey = "alt"   // Synonym for ModAlt
	ModCtrl  ModKey = "ctrl"  // Alternate action for ^↩
	ModShift ModKey = "shift" // Alternate action for ⇧↩
	ModFn    ModKey = "fn"    // Alternate action for fn↩
)

// Mod element gives you control over how the modifier keys react
type Mod struct {
	Variables map[string]string `json:"variables,omitempty"`
	Valid     *bool             `json:"valid,omitempty"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Icon      *Icon             `json:"icon,omitempty"`
}

// Text element defines the text the user will get when copying the selected result row
type Text struct {
	Copy      string `json:"copy,omitempty"`
	Largetype string `json:"largetype,omitempty"`
}

// Item a workflow object
type Item struct {
	Variables    map[string]string `json:"variables,omitempty"`
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Icon         *Icon             `json:"icon,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Type         string            `json:"type,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Match        string            `json:"match,omitempty"`
	Mods         map[ModKey]Mod    `json:"mods,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	QuicklookURL string            `json:"quicklookurl,omitempty"`
}

// Rerun re-run automatically after an interval
type Rerun float64

// Variables passed out of the script filter within a `variables` object.
type Variables map[string]string

// Items array of `item`
type Items []*Item

// ScriptFilter JSON Format
type ScriptFilter struct {
	rerun     Rerun
	variables Variables
	items     Items
}

type out struct {
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

const fatalErrorJSON = `{"items": [{"title": "Fatal Error","subtitle": "%s",}]}`

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() ScriptFilter {
	return ScriptFilter{}
}

// Append a new Item to Items
func (s *ScriptFilter) Append(item *Item) {
	s.items = append(s.items, item)
}

// Marshal ScriptFilter as Json
func (s *ScriptFilter) Marshal() []byte {
	res, err := json.Marshal(
		out{
			Rerun:     s.rerun,
			Variables: s.variables,
			Items:     s.items,
		},
	)
	if err != nil {
		return []byte(fmt.Sprintf(fatalErrorJSON, err.Error()))
	}

	return res
}

// Workflow is map of ScriptFilters
type Workflow struct {
	std     ScriptFilter
	warn    ScriptFilter
	err     ScriptFilter
	streams streams
}

type streams struct {
	out io.Writer
	err io.Writer
}

// SetStreams redirect stdout and stderr to s
func (w *Workflow) SetStreams(out, err io.Writer) {
	w.streams.out = out
	w.streams.err = err
}

// NewWorkflow has simple ScriptFilter api
func NewWorkflow() *Workflow {
	return &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
			err: os.Stdout,
		},
	}
}

// Append a new Item to standard ScriptFilter
func (w *Workflow) Append(item *Item) {
	w.std.Append(item)
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
	if len(w.std.items) == 0 {
		return w.warn.Marshal()
	}

	return w.std.Marshal()
}

// Fatal output error to io stream
func (w *Workflow) Fatal(title, subtitle string) {
	w.error(title, subtitle)
	res := w.err.Marshal()
	fmt.Fprintln(w.streams.err, string(res))
}

// Output to io stream
func (w *Workflow) Output() {
	res := w.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
}
