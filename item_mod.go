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

// Mod element gives you control over how the modifier keys react
type Mod struct {
	variables Variables
	valid     bool
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
	m.valid = b
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

type iMod struct {
	Variables Variables `json:"variables,omitempty"`
	Valid     bool      `json:"valid,omitempty"`
	Arg       string    `json:"arg,omitempty"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Icon      *Icon     `json:"icon,omitempty"`
}

func (m *Mod) MarshalJSON() ([]byte, error) {
	out := &iMod{
		Variables: m.variables,
		Valid:     m.valid,
		Arg:       m.arg,
		Subtitle:  m.subtitle,
		Icon:      m.icon,
	}
	return json.Marshal(out)
}

func (m *Mod) UnmarshalJSON(data []byte) error {
	in := &iMod{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*m = Mod{
		variables: in.Variables,
		valid:     in.Valid,
		arg:       in.Arg,
		subtitle:  in.Subtitle,
		icon:      in.Icon,
	}
	return nil
}
