package alfred

import (
	"reflect"

	"github.com/sahilm/fuzzy"
)

type ItemProperty int

const (
	ItemPropertyTitle ItemProperty = iota
	ItemPropertySubtitle
	ItemPropertyArg
	ItemPropertyUID
)

func (p ItemProperty) String() string {
	var field string
	switch p {
	case ItemPropertyTitle:
		field = "title"
	case ItemPropertySubtitle:
		field = "subtitle"
	case ItemPropertyArg:
		field = "arg"
	case ItemPropertyUID:
		field = "uid"
	}
	return field
}

// Filter by item title with fuzzy
func (w *Workflow) Filter(query string) *Workflow {
	w.std.items = w.std.items.Filter(query)
	return w
}

func (w *Workflow) FilterByItemProperty(f func(s string) bool, property ItemProperty) *Workflow {
	items := make(Items, 0, cap(w.std.items))
	for _, item := range w.std.items {
		v := getItemValue(item, property.String())
		if f(v) {
			items = append(items, item)
		}
	}

	// update
	w.std.items = items
	return w
}

func getItemValue(item *Item, field string) string {
	rv := reflect.Indirect(reflect.ValueOf(item))
	rt := rv.Type()

	f, ok := rt.FieldByName(field)
	if !ok {
		return ""
	}

	v := rv.FieldByName(f.Name).String()
	return v
}

// String retruns a title of Item for fuzzy interface
func (i Items) String(idx int) string {
	return i[idx].title
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
