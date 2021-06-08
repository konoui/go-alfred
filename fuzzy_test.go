package alfred

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name  string
		input Items
		args  args
		want  Items
	}{
		{
			name:  "return all items if empty query",
			input: items01,
			want:  items01,
			args: args{
				query: "",
			},
		},
		{
			name:  "perfect matching",
			input: items01,
			want: Items{
				items01[0],
			},
			args: args{
				query: items01[0].title,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Filter(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflow_FilterByItemProperty(t *testing.T) {
	type args struct {
		f        func(s string) bool
		property FilterProperty
	}
	tests := []struct {
		name  string
		args  args
		items Items
		want  Items
	}{
		{
			name: "by title",
			args: args{
				f:        func(s string) bool { return strings.HasPrefix(s, "title2") },
				property: FilterTitle,
			},
			items: items05,
			want: Items{
				items05[1],
			},
		},
		{
			name: "by subtitle",
			args: args{
				f:        func(s string) bool { return strings.HasPrefix(s, "subtitle2") },
				property: FilterSubtitle,
			},
			items: items05,
			want: Items{
				items05[1],
			},
		},
		{
			name: "by arg",
			args: args{
				f:        func(s string) bool { return strings.HasPrefix(s, "arg2") },
				property: FilterArg,
			},
			items: items05,
			want: Items{
				items05[1],
			},
		},
		{
			name: "by uid",
			args: args{
				f:        func(s string) bool { return strings.HasPrefix(s, "uid2") },
				property: FilterUID,
			},
			items: items05,
			want: Items{
				items05[1],
			},
		},
		{
			name: "not match",
			args: args{
				f:        func(s string) bool { return strings.HasPrefix(s, "uid100") },
				property: FilterUID,
			},
			items: items05,
			want:  Items{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow().Append(tt.items...)
			got := w.FilterByItemProperty(tt.args.f, tt.args.property).std.items
			if diff := Diff(tt.want, got); diff != "" {
				t.Errorf("-want +got\n%s", diff)
			}
		})
	}
}
