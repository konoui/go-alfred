package alfred

import (
	"github.com/sahilm/fuzzy"
)

// String retrun a title of Item for fuzzy interface
func (i Items) String(idx int) string {
	return i[idx].Title
}

// Len return length of Items for fuzzy interface
func (i Items) Len() int {
	return len(i)
}

// Filter fuzzy search items using query
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
