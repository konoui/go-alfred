package alfred

import (
	"reflect"
	"testing"
)

var filterItems = Items{
	&Item{
		Title: "aaaaaa",
	},
	&Item{
		Title: "bbbbbb",
	},
}

func TestFilter(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		i    Items
		args args
		want Items
	}{
		{
			name: "all items if empty query",
			i:    filterItems,
			want: filterItems,
			args: args{
				query: "",
			},
		},
		{
			name: "Perfect matching",
			i:    filterItems,
			want: Items{
				filterItems[0],
			},
			args: args{
				query: "aaaaaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Filter(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
