package alfred

import (
	"github.com/sahilm/fuzzy"
)

// Filter searches from current items
func (w *Workflow) Filter(query string) *Workflow {
	w.std.Items = w.std.Items.Filter(query)
	return w
}

// String retruns a title of Item for fuzzy interface
func (i Items) String(idx int) string {
	return i[idx].Title
}

// Len returns length of Items for fuzzy interface
func (i Items) Len() int {
	return len(i)
}

// Filter searches items using query
func (i Items) Filter(query string) Items {
	if query == "" {
		return i
	}
	results := fuzzy.FindFrom(query, i)
	items := make(Items, results.Len())
	for idx, r := range results {
		items[idx] = i[r.Index]
	}

	return items
}
