package alfred

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// DiffScriptFilter is a helper function that compare unsorted ScriptFilter Output
// return "" if `gotData` is equal to `wantData` regardless of sorted or unsorted
func DiffScriptFilter(wantData, gotData []byte) string {
	want := &ScriptFilter{}
	got := &ScriptFilter{}
	if err := json.Unmarshal(wantData, want); err != nil {
		return fmt.Sprintf("Unmarshal Error in wantData: %+v\n, string(wantData): %s\n, string(gotData): %s\n", err, string(wantData), string(gotData))
	}

	if err := json.Unmarshal(gotData, got); err != nil {
		return fmt.Sprintf("Unmarshal Error in gotData: %+v\n, string(wantData): %s\n, string(gotData): %s\n", err, string(wantData), string(gotData))
	}

	sort.Slice(want.Items, func(i, j int) bool {
		return want.Items[i].Title < want.Items[j].Title
	})

	sort.Slice(got.Items, func(i, j int) bool {
		return got.Items[i].Title < got.Items[j].Title
	})

	return cmp.Diff(want, got)
}
