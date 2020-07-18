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
