package alfred

import (
	"errors"
	"testing"
	"time"
)

func TestWorkflow_Cache(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		wf   *Workflow
		args args
	}{
		{
			name: "Cache behaves singleton. return same address",
			wf:   NewWorkflow(),
			args: args{
				key: "test1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.wf.Cache(tt.args.key)
			got := tt.wf.Cache(tt.args.key)
			if err := tt.wf.Cache(tt.args.key).Err(); err != nil {
				t.Errorf("unexpected error %s", err)
			}
			if want != got {
				t.Errorf("Cache() does not behave singleton want %#v\n, got %#v\n", want, got)
			}
		})
	}
}

func TestCache_Workflow(t *testing.T) {
	tests := []struct {
		name string
		want *Workflow
	}{
		{
			name: "workflow.Cache(key).Wrokflow() return itself",
			want: NewWorkflow(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.want.Cache("dummy").Workflow(); got != tt.want {
				t.Errorf("Cache.Workflow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_LoadStoreClearItems(t *testing.T) {
	tests := []struct {
		name      string
		wf        *Workflow
		ttl       time.Duration
		expectErr bool
	}{
		{
			name:      "equal address and data stored and loaded instance",
			wf:        NewWorkflow(),
			ttl:       1 * time.Minute,
			expectErr: false,
		},
		{
			name:      "cache expired error",
			wf:        NewWorkflow(),
			ttl:       0 * time.Minute,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheKey := "test1"
			defer func() {
				if err := tt.wf.Cache(cacheKey).ClearItems().Err(); err != nil {
					t.Error(err)
				}
			}()
			// Input test data
			prepared := NewWorkflow().Append(items01[0])
			err := prepared.Cache(cacheKey).StoreItems().Err()
			if err != nil {
				t.Fatal(err)
			}

			// load cache from new workflow
			err = tt.wf.Cache(cacheKey).TTL(tt.ttl).LoadItems().Err()
			if !tt.expectErr && err != nil {
				t.Fatal(err)
			}

			if tt.expectErr && !errors.Is(err, ErrCacheExpired) {
				t.Errorf("want: expired cache error, got: %#v\n", err)
			}

			// compare new workflow data to soted workflow data
			want := prepared.std.items
			got := tt.wf.std.items
			if diff := Diff(want, got); !tt.expectErr && diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func Test_SetGetCacheSuffix(t *testing.T) {
	name := "set/get cache suffix"
	t.Run(name, func(t *testing.T) {
		want := "test"
		awf := NewWorkflow(
			WithCacheSuffix(want),
		)
		got := awf.getCacheSuffix()
		if want != got {
			t.Errorf("want %s got %s", want, got)
		}
	})
}
