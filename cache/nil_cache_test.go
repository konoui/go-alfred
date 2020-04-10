package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestNewNilCache(t *testing.T) {
	tests := []struct {
		name string
		want *NilCache
	}{
		{
			name: "valid directory",
			want: &NilCache{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNilCache()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %+v\ngot: %+v", tt.want, got)
			}
		})
	}
}

func TestNilCache(t *testing.T) {
	tests := []struct {
		name          string
		expiredResult bool
		loadResult    error
		storeResult   error
	}{
		{
			name:          "cacher interface",
			expiredResult: true,
			loadResult:    nil,
			storeResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewNilCache()
			if got := c.Expired(90 * time.Second); got != tt.expiredResult {
				t.Errorf("want: %+v\ngot: %+v", tt.expiredResult, got)
			}

			if err := c.Clear(); err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			if got := c.Load(tt.storeResult); got != tt.loadResult {
				t.Errorf("want: %+v\ngot: %+v", tt.loadResult, got)
			}

			if got := c.Store(tt.storeResult); got != tt.storeResult {
				t.Errorf("want: %+v\ngot: %+v", tt.storeResult, got)
			}
		})
	}
}
