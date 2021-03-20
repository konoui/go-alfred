package alfred

import "encoding/json"

// ModKey is a mod key pressed by the user to run an alternate
type ModKey string

// Valid attribute to mark if the result is valid based on the modifier selection and Add a different arg to be passed out if actioned with the modifier.
const (
	ModCmd   ModKey = "cmd"   // Alternate action for ⌘↩
	ModAlt   ModKey = "alt"   // Alternate action for ⌥↩
	ModOpt   ModKey = "alt"   // Synonym for ModAlt
	ModCtrl  ModKey = "ctrl"  // Alternate action for ^↩
	ModShift ModKey = "shift" // Alternate action for ⇧↩
	ModFn    ModKey = "fn"    // Alternate action for fn↩
)

type Mods map[ModKey]*Mod

// Mod element gives you control over how the modifier keys react
type Mod struct {
	variables Variables
	valid     *bool
	arg       string
	subtitle  string
	icon      *Icon
}

// NewMod generates new mod
func NewMod() *Mod {
	return new(Mod)
}

// Variables adds mod variables
func (m *Mod) Variables(vars Variables) *Mod {
	for k, v := range vars {
		m.Variable(k, v)
	}
	return m
}

// Variable adds mod single key/value pair
func (m *Mod) Variable(k, v string) *Mod {
	if m.variables == nil {
		m.variables = make(Variables)
	}
	m.variables[k] = v
	return m
}

// Valid adds mod valid or not
func (m *Mod) Valid(b bool) *Mod {
	m.valid = &b
	return m
}

// Arg adds mod argument
func (m *Mod) Arg(s string) *Mod {
	m.arg = s
	return m
}

// Subtitle adds mod subtitle
func (m *Mod) Subtitle(s string) *Mod {
	m.subtitle = s
	return m
}

// Icon adds mod icon
func (m *Mod) Icon(icon *Icon) *Mod {
	m.icon = icon
	return m
}

func (m *Mod) MarshalJSON() ([]byte, error) {
	out := m.internal()
	return json.Marshal(out)
}

func (m *Mod) UnmarshalJSON(data []byte) error {
	in := &iMod{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*m = *in.external()
	return nil
}

type iMods map[ModKey]*iMod

type iMod struct {
	Variables Variables `json:"variables,omitempty"`
	Valid     *bool     `json:"valid,omitempty"`
	Arg       string    `json:"arg,omitempty"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Icon      *iIcon    `json:"icon,omitempty"`
}

func (m Mods) internal() iMods {
	mods := make(map[ModKey]*iMod)
	for k, v := range m {
		mods[k] = v.internal()
	}
	return mods
}

func (m iMods) external() Mods {
	mods := make(map[ModKey]*Mod)
	for k, v := range m {
		mods[k] = v.external()
	}
	return mods
}

func (m *Mod) internal() *iMod {
	if m == nil {
		return nil
	}
	return &iMod{
		Variables: m.variables,
		Valid:     m.valid,
		Arg:       m.arg,
		Subtitle:  m.subtitle,
		Icon:      m.icon.internal(),
	}
}

func (m *iMod) external() *Mod {
	if m == nil {
		return nil
	}
	return &Mod{
		variables: m.Variables,
		valid:     m.Valid,
		arg:       m.Arg,
		subtitle:  m.Subtitle,
		icon:      m.Icon.external(),
	}
}
