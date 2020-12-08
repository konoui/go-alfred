package alfred

import "encoding/json"

// Icon displayed in the result row
type Icon struct {
	typ  string
	path string
}

// NewIcon generates new icon
func NewIcon() *Icon {
	return new(Icon)
}

// Type adds icon type
func (i *Icon) Type(s string) *Icon {
	i.typ = s
	return i
}

// Path adds icon abs path
func (i *Icon) Path(s string) *Icon {
	i.path = s
	return i
}

func (i *Icon) MarshalJSON() ([]byte, error) {
	out := i.internal()
	return json.Marshal(out)
}

func (i *Icon) UnmarshalJSON(data []byte) error {
	in := &iIcon{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*i = *in.external()
	return nil
}

type iIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

func (i *Icon) internal() *iIcon {
	if i == nil {
		return nil
	}
	return &iIcon{
		Type: i.typ,
		Path: i.path,
	}
}

func (i *iIcon) external() *Icon {
	if i == nil {
		return nil
	}
	return &Icon{
		typ:  i.Type,
		path: i.Path,
	}
}
