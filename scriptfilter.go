package alfred

import (
	"encoding/json"
	"fmt"
)

const (
	fatalErrorJSON = `{"items": [{"title": "Fatal Error","subtitle": "%s",}]}`
	sentMessage    = "The workflow has already sent"
)

// ScriptFilter JSON Format
type ScriptFilter struct {
	rerun     Rerun
	variables Variables
	items     Items
}

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() *ScriptFilter {
	return &ScriptFilter{}
}

// Variables sets ScriptFilter variables
func (s *ScriptFilter) Variables(vars Variables) {
	for k, v := range vars {
		s.Variable(k, v)
	}
}

// Variable sets ScriptFilter variable
func (s *ScriptFilter) Variable(k, v string) {
	if s.variables == nil {
		s.variables = make(Variables)
	}
	s.variables[k] = v
}

// Rerun adds rerun variable
func (s *ScriptFilter) Rerun(i Rerun) {
	s.rerun = i
}

// Append adds item
func (s *ScriptFilter) Append(i ...*Item) {
	s.items = append(s.items, i...)
}

// Clear remove all items
func (s *ScriptFilter) Clear() {
	s.items = Items{}
}

// Marshal ScriptFilter as Json
func (s *ScriptFilter) Marshal() []byte {
	res, err := json.Marshal(s)
	if err != nil {
		return []byte(fmt.Sprintf(fatalErrorJSON, err.Error()))
	}

	return res
}

// IsEmpty return true if the items is empty
func (s *ScriptFilter) IsEmpty() bool {
	return len(s.items) == 0
}

type iScriptFilter struct {
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

func (s *ScriptFilter) MarshalJSON() ([]byte, error) {
	out := &iScriptFilter{
		Rerun:     s.rerun,
		Variables: s.variables,
		Items:     s.items,
	}
	return json.Marshal(out)
}

func (s *ScriptFilter) UnmarshalJSON(data []byte) error {
	in := &iScriptFilter{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*s = ScriptFilter{
		rerun:     in.Rerun,
		variables: in.Variables,
		items:     in.Items,
	}
	return nil
}
