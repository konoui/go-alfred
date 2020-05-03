package alfred

import (
	"io"
	"sync"
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
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

// Workflow is map of ScriptFilters
type Workflow struct {
	std     ScriptFilter
	warn    ScriptFilter
	err     ScriptFilter
	caches  sync.Map
	streams streams
}

type streams struct {
	out io.Writer
}
