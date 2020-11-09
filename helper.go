package alfred

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// DiffScriptFilter is a helper function that compare unsorted ScriptFilter Output
// return "" if `gotData` is equal to `wantData` regardless of sorted or unsorted
func DiffScriptFilter(wantData, gotData []byte) string {
	var out1, out2 ScriptFilter
	if err := decodeObjects(wantData, gotData, &out1, &out2); err != nil {
		return err.Error()
	}

	sort.Slice(out1.items, func(i, j int) bool {
		return out1.items[i].title < out1.items[j].title
	})

	sort.Slice(out2.items, func(i, j int) bool {
		return out2.items[i].title < out2.items[j].title
	})
	return diffStringByEncode(&out1, &out2)
}

func diffStringByEncode(want, got interface{}) string {
	out1Data, out2Data, err := encodeObjects(want, got)
	if err != nil {
		return err.Error()
	}
	return cmp.Diff(string(out1Data), string(out2Data))
}

func encodeObjects(want, got interface{}) (out1, out2 []byte, err error) {
	out1, err = json.Marshal(want)
	if err != nil {
		return nil, nil, fmt.Errorf("1st argument is unable to encode to json due to %s", err.Error())
	}

	out2, err = json.Marshal(got)
	if err != nil {
		return nil, nil, fmt.Errorf("2nd argument is unable to encode to json due to %s", err.Error())
	}
	return
}

func decodeObjects(wantData, gotData []byte, wantOut, gotOut interface{}) error {
	errA := json.NewDecoder(bytes.NewReader(wantData)).Decode(wantOut)
	errB := json.NewDecoder(bytes.NewReader(gotData)).Decode(gotOut)
	if errA != nil && errB != nil {
		return errors.New("both arguments are invalid json")
	}
	if errA != nil {
		return errors.New("1st argument is invalid json")
	}
	if errB != nil {
		return errors.New("2nd argument is invalid json")
	}
	return nil
}
