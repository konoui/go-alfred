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
