package alfred

import (
	"encoding/json"
)

// Text element defines the text the user will get when copying the selected result row
type Text struct {
	copy      string
	largeType string
}

// NewText generates new text
func NewText() *Text {
	return new(Text)
}

// Copy adds text copy value
func (t *Text) Copy(s string) *Text {
	t.copy = s
	return t
}

// LargeType adds text large type
func (t *Text) LargeType(s string) *Text {
	t.largeType = s
	return t
}

type iText struct {
	Copy      string `json:"copy,omitempty"`
	LargeType string `json:"largetype,omitempty"`
}

func (t *Text) MarshalJSON() ([]byte, error) {
	out := &iText{
		Copy:      t.copy,
		LargeType: t.largeType,
	}
	return json.Marshal(out)
}

func (t *Text) UnmarshalJSON(data []byte) error {
	in := &iText{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*t = Text{
		copy:      in.Copy,
		largeType: in.LargeType,
	}
	return nil
}
