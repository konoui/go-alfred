package alfred

import (
	"reflect"
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
			input: item01,
			want:  item01,
			args: args{
				query: "",
			},
		},
		{
			name:  "perfect matching",
			input: item01,
			want: Items{
				item01[0],
			},
			args: args{
				query: item01[0].Title,
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
