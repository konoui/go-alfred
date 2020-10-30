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
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() *ScriptFilter {
	return &ScriptFilter{}
}

// Append a new Item to Items
func (s *ScriptFilter) Append(item *Item) {
	s.Items = append(s.Items, item)
}

// SetVariables sets ScriptFilter variables
func (s *ScriptFilter) SetVariables(vars Variables) {
	for k, v := range vars {
		s.SetVariable(k, v)
	}
}

// SetVariable sets ScriptFilter variable
func (s *ScriptFilter) SetVariable(k, v string) {
	if s.Variables == nil {
		s.Variables = make(Variables)
	}
	s.Variables[k] = v
}

// SetRerun sets rerun variable
func (s *ScriptFilter) SetRerun(i Rerun) {
	s.Rerun = i
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
	return len(s.Items) == 0
}
