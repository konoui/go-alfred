package alfred

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// DiffOutput is a helper function that compare unsorted ScriptFilter Output
// return "" if `gotData` is equal to `wantData` regardless of sorted or unsorted
func DiffOutput(wantData, gotData []byte) string {
	out1, out2 := new(ScriptFilter), new(ScriptFilter)
	errA := json.NewDecoder(bytes.NewReader(wantData)).Decode(out1)
	errB := json.NewDecoder(bytes.NewReader(gotData)).Decode(out2)
	if errA != nil {
		return "wantData is invalid json"
	}
	if errB != nil {
		return "gotData is invalid json"
	}

	return Diff(out1, out2)
}

func Diff(want, got interface{}) string {
	wrt := reflect.TypeOf(want)
	grt := reflect.TypeOf(got)
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		return fmt.Sprintf("want is %v but got is %v", wrt, grt)
	}

	switch val := want.(type) {
	case *Mod:
		w := want.(*Mod).internal()
		g := got.(*Mod).internal()
		return cmp.Diff(w, g)
	case *Text:
		w := want.(*Text).internal()
		g := got.(*Text).internal()
		return cmp.Diff(w, g)
	case *Icon:
		w := want.(*Icon).internal()
		g := got.(*Icon).internal()
		return cmp.Diff(w, g)
	case *Item:
		w := want.(*Item).internal()
		g := got.(*Item).internal()
		return cmp.Diff(w, g)
	case Mods:
		w := want.(Mods).internal()
		g := got.(Mods).internal()
		return cmp.Diff(w, g)
	case Items:
		w := want.(Items).internal()
		g := got.(Items).internal()
		sortItems(w, g)
		return cmp.Diff(w, g)
	case *ScriptFilter:
		w := want.(*ScriptFilter).internal()
		g := got.(*ScriptFilter).internal()
		sortItems(w.Items, g.Items)
		return cmp.Diff(w, g)
	default:
		return fmt.Sprintf("unsupported type in want/got %v", val)
	}
}

func sortItems(in1, in2 iItems) {
	sort.Slice(in1, func(i, j int) bool {
		return in1[i].Title < in1[j].Title
	})

	sort.Slice(in2, func(i, j int) bool {
		return in2[i].Title < in2[j].Title
	})
}
