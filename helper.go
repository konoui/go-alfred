package alfred

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// DiffScriptFilter is a helper function that compare unsorted ScriptFilter Output
// return "" if `gotData` is equal to `wantData` regardless of sorted or unsorted
func DiffScriptFilter(wantData, gotData []byte) string {
	out1, out2 := new(iScriptFilter), new(iScriptFilter)
	errA := json.NewDecoder(bytes.NewReader(wantData)).Decode(out1)
	errB := json.NewDecoder(bytes.NewReader(gotData)).Decode(out2)
	if errA != nil {
		return "wantData is invalid json"
	}
	if errB != nil {
		return "gotData is invalid json"
	}

	sort.Slice(out1.Items, func(i, j int) bool {
		return out1.Items[i].title < out1.Items[j].title
	})

	sort.Slice(out2.Items, func(i, j int) bool {
		return out2.Items[i].title < out2.Items[j].title
	})

	out1Data, err := json.Marshal(out1)
	if err != nil {
		return fmt.Sprintf("failed to marshal wantData due to %v", err)
	}

	out2Data, err := json.Marshal(out2)
	if err != nil {
		return fmt.Sprintf("failed to marshal gotData due to %v", err)
	}

	return cmp.Diff(string(out1Data), string(out2Data))
}
