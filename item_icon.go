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

type iIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

func (i *Icon) MarshalJSON() ([]byte, error) {
	out := &iIcon{
		Type: i.typ,
		Path: i.path,
	}
	return json.Marshal(out)
}

func (i *Icon) UnmarshalJSON(data []byte) error {
	in := &iIcon{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	*i = Icon{
		typ:  in.Type,
		path: in.Path,
	}
	return nil
}
