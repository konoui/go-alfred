// Package alfred defines script filter structures and provides simple apis
// see https://www.alfredapp.com/help/workflows/inputs/script-filter/json/
package alfred

// Rerun re-run automatically after an interval
type Rerun float64

// Variables passed out of the script filter within a `variables` object.
type Variables map[string]string

// Items array of `item`
type Items []*Item

// Item a workflow object
type Item struct {
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

// NewItem generates new item
func NewItem() *Item {
	return new(Item)
}

// SetVariables sets variables
func (i *Item) SetVariables(vars Variables) *Item {
	for k, v := range vars {
		i.SetVariable(k, v)
	}
	return i
}

// SetVariable sets single key/value pair
func (i *Item) SetVariable(k, v string) *Item {
	if i.Variables == nil {
		i.Variables = make(Variables)
	}
	i.Variables[k] = v
	return i
}

// SetUID sets uid
func (i *Item) SetUID(s string) *Item {
	i.UID = s
	return i
}

// SetTitle sets title
func (i *Item) SetTitle(s string) *Item {
	i.Title = s
	return i
}

// SetSubtitle sets subtitle
func (i *Item) SetSubtitle(s string) *Item {
	i.Subtitle = s
	return i
}

// SetArg sets arg
func (i *Item) SetArg(arg string) *Item {
	i.Arg = arg
	return i
}

// SetIcon sets icon
func (i *Item) SetIcon(icon *Icon) *Item {
	i.Icon = icon
	return i
}

// SetAutocomplete sets autocomplete
func (i *Item) SetAutocomplete(s string) *Item {
	i.Autocomplete = s
	return i
}

// SetValid sets valid
func (i *Item) SetValid(b bool) *Item {
	i.Valid = b
	return i
}

// SetMatch sets match
func (i *Item) SetMatch(s string) *Item {
	i.Match = s
	return i
}

// SetMods sets mods
func (i *Item) SetMods(mods map[ModKey]*Mod) *Item {
	for k, v := range mods {
		i.SetMod(k, v)
	}
	return i
}

// SetMod sets single mod
func (i *Item) SetMod(key ModKey, mod *Mod) *Item {
	if i.Mods == nil {
		i.Mods = make(map[ModKey]*Mod)
	}
	i.Mods[key] = mod
	return i
}

// SetText sets text
func (i *Item) SetText(t *Text) *Item {
	i.Text = t
	return i
}

// SetQuicklookURL sets quicklookURL
func (i *Item) SetQuicklookURL(u string) *Item {
	i.QuicklookURL = u
	return i
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
	Variables Variables `json:"variables,omitempty"`
	Valid     bool      `json:"valid,omitempty"`
	Arg       string    `json:"arg,omitempty"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Icon      *Icon     `json:"icon,omitempty"`
}

// NewMod generates new mod
func NewMod() *Mod {
	return new(Mod)
}

// SetVariables sets mod variables
func (m *Mod) SetVariables(vars Variables) *Mod {
	for k, v := range vars {
		m.SetVariable(k, v)
	}
	return m
}

// SetVariable sets mod single key/value pair
func (m *Mod) SetVariable(k, v string) *Mod {
	if m.Variables == nil {
		m.Variables = make(Variables)
	}
	m.Variables[k] = v
	return m
}

// SetValid sets mod valid or not
func (m *Mod) SetValid(b bool) *Mod {
	m.Valid = b
	return m
}

// SetArg sets mod argument
func (m *Mod) SetArg(s string) *Mod {
	m.Arg = s
	return m
}

// SetSubtitle sets mod subtitle
func (m *Mod) SetSubtitle(s string) *Mod {
	m.Subtitle = s
	return m
}

// SetIcon sets mod icon
func (m *Mod) SetIcon(icon *Icon) *Mod {
	m.Icon = icon
	return m
}

// Icon displayed in the result row
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// NewIcon generates new icon
func NewIcon() *Icon {
	return new(Icon)
}

// SetType sets icon type
func (i *Icon) SetType(s string) *Icon {
	i.Type = s
	return i
}

// SetPath sets icon abs path
func (i *Icon) SetPath(s string) *Icon {
	i.Path = s
	return i
}

// Text element defines the text the user will get when copying the selected result row
type Text struct {
	Copy      string `json:"copy,omitempty"`
	LargeType string `json:"largetype,omitempty"`
}

// NewText generates new text
func NewText() *Text {
	return new(Text)
}

// SetCopy sets text copy value
func (t *Text) SetCopy(s string) *Text {
	t.Copy = s
	return t
}

// SetLargeType sets text large type
func (t *Text) SetLargeType(s string) *Text {
	t.LargeType = s
	return t
}
