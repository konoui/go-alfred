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
	valid        *bool
	match        string
	mods         Mods
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
	i.valid = &b
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

func (i *Item) MarshalJSON() ([]byte, error) {
	out := i.internal()
	return json.Marshal(out)
}

func (i *Item) UnmarshalJSON(data []byte) error {
	in := &iItem{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*i = *in.external()
	return nil
}

type iItems []*iItem

type iItem struct {
	Variables    Variables `json:"variables,omitempty"`
	UID          string    `json:"uid,omitempty"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle,omitempty"`
	Arg          string    `json:"arg,omitempty"`
	Icon         *iIcon    `json:"icon,omitempty"`
	Autocomplete string    `json:"autocomplete,omitempty"`
	Type         string    `json:"type,omitempty"`
	Valid        *bool     `json:"valid,omitempty"`
	Match        string    `json:"match,omitempty"`
	Mods         iMods     `json:"mods,omitempty"`
	Text         *iText    `json:"text,omitempty"`
	QuicklookURL string    `json:"quicklookurl,omitempty"`
}

func (i Items) internal() iItems {
	items := make(iItems, len(i), cap(i))
	for idx, itm := range i {
		items[idx] = itm.internal()
	}
	return items
}

func (i iItems) external() Items {
	items := make(Items, len(i), cap(i))
	for idx, itm := range i {
		items[idx] = itm.external()
	}
	return items
}

func (i *Item) internal() *iItem {
	return &iItem{
		Variables:    i.variables,
		UID:          i.uid,
		Title:        i.title,
		Subtitle:     i.subtitle,
		Arg:          i.arg,
		Icon:         i.icon.internal(),
		Autocomplete: i.autocomplete,
		Type:         i.typ,
		Valid:        i.valid,
		Match:        i.match,
		Mods:         i.mods.internal(),
		Text:         i.text.internal(),
		QuicklookURL: i.quicklookURL,
	}
}

func (i *iItem) external() *Item {
	return &Item{
		variables:    i.Variables,
		uid:          i.UID,
		title:        i.Title,
		subtitle:     i.Subtitle,
		arg:          i.Arg,
		icon:         i.Icon.external(),
		autocomplete: i.Autocomplete,
		typ:          i.Type,
		valid:        i.Valid,
		match:        i.Match,
		mods:         i.Mods.external(),
		text:         i.Text.external(),
		quicklookURL: i.QuicklookURL,
	}
}
