// Package alfred defines script filter structures and provides simple apis
// see https://www.alfredapp.com/help/workflows/inputs/script-filter/json/
package alfred

import (
	"encoding/json"
)

// Rerun re-run automatically after an interval
type Rerun float64

// Variables passed out of the script filter within a `variables` object.
type Variables map[string]string

// Items array of `item`
type Items []*Item

// Item a workflow object
type Item struct {
	variables    Variables
	uid          string
	title        string
	subtitle     string
	arg          string
	icon         *Icon
	autocomplete string
	typ          string
	valid        bool
	match        string
	mods         map[ModKey]*Mod
	text         *Text
	quicklookURL string
}

// NewItem generates new item
func NewItem() *Item {
	return new(Item)
}

// Variables adds variables
func (i *Item) Variables(vars Variables) *Item {
	for k, v := range vars {
		i.Variable(k, v)
	}
	return i
}

// Variable adds single key/value pair
func (i *Item) Variable(k, v string) *Item {
	if i.variables == nil {
		i.variables = make(Variables)
	}
	i.variables[k] = v
	return i
}

// UID adds uid
func (i *Item) UID(s string) *Item {
	i.uid = s
	return i
}

// Title adds title
func (i *Item) Title(s string) *Item {
	i.title = s
	return i
}

// Subtitle adds subtitle
func (i *Item) Subtitle(s string) *Item {
	i.subtitle = s
	return i
}

// Arg adds arg
func (i *Item) Arg(arg string) *Item {
	i.arg = arg
	return i
}

// Icon adds icon
func (i *Item) Icon(icon *Icon) *Item {
	i.icon = icon
	return i
}

// Autocomplete adds autocomplete
func (i *Item) Autocomplete(s string) *Item {
	i.autocomplete = s
	return i
}

// Valid adds valid
func (i *Item) Valid(b bool) *Item {
	i.valid = b
	return i
}

// Match adds match
func (i *Item) Match(s string) *Item {
	i.match = s
	return i
}

// Mods adds mods
func (i *Item) Mods(mods map[ModKey]*Mod) *Item {
	for k, v := range mods {
		i.Mod(k, v)
	}
	return i
}

// Mod adds single mod
func (i *Item) Mod(key ModKey, mod *Mod) *Item {
	if i.mods == nil {
		i.mods = make(map[ModKey]*Mod)
	}
	i.mods[key] = mod
	return i
}

// Text adds text
func (i *Item) Text(t *Text) *Item {
	i.text = t
	return i
}

// QuicklookURL adds quicklookURL
func (i *Item) QuicklookURL(u string) *Item {
	i.quicklookURL = u
	return i
}

type iItem struct {
	Variables    Variables       `json:"variables,omitempty"`
	UID          string          `json:"uid,omitempty"`
	Title        string          `json:"title"`
	Subtitle     string          `json:"subtitle,omitempty"`
	Arg          string          `json:"arg,omitempty"`
	Icon         *Icon           `json:"icon,omitempty"`
	Autocomplete string          `json:"autocomplete,omitempty"`
	Type         string          `json:"type,omitempty"`
	Valid        bool            `json:"valid,omitempty"`
	Match        string          `json:"match,omitempty"`
	Mods         map[ModKey]*Mod `json:"mods,omitempty"`
	Text         *Text           `json:"text,omitempty"`
	QuicklookURL string          `json:"quicklookurl,omitempty"`
}

func (i *Item) MarshalJSON() ([]byte, error) {
	out := &iItem{
		Variables:    i.variables,
		UID:          i.uid,
		Title:        i.title,
		Subtitle:     i.subtitle,
		Arg:          i.arg,
		Icon:         i.icon,
		Autocomplete: i.autocomplete,
		Type:         i.typ,
		Valid:        i.valid,
		Match:        i.match,
		Mods:         i.mods,
		Text:         i.text,
		QuicklookURL: i.quicklookURL,
	}
	return json.Marshal(out)
}

func (i *Item) UnmarshalJSON(data []byte) error {
	in := &iItem{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*i = Item{
		variables:    in.Variables,
		uid:          in.UID,
		title:        in.Title,
		subtitle:     in.Subtitle,
		arg:          in.Arg,
		icon:         in.Icon,
		autocomplete: in.Autocomplete,
		typ:          in.Type,
		valid:        in.Valid,
		match:        in.Match,
		mods:         in.Mods,
		text:         in.Text,
		quicklookURL: in.QuicklookURL,
	}
	return nil
}
