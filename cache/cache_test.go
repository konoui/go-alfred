package cache

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var tmpDir = os.TempDir()

type example struct {
	A string
	B string
	C []string
}

var storedValue = example{
	A: "AAAAA",
	B: "BBBBBB",
	C: []string{
		"11111",
		"22222",
		"33333",
	},
}

func TestNewCache(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		file      string
		expectErr bool
	}{
		{
			name:      "valid directory",
			dir:       tmpDir,
			file:      "test1",
			expectErr: false,
		},
		{
			name:      "invalid directory",
			dir:       "/unk",
			file:      "test2",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.dir, tt.file)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %+v", err)
			}
		})
	}
}

func TestStore(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		file      string
		expectErr bool
	}{
		{
			name:      "create cache file on temp dir",
			dir:       tmpDir,
			file:      "test1",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := New(tt.dir, tt.file)
			if err != nil {
				t.Fatal(err)
			}

			// remove cache file before test
			if err = cache.Clear(); err != nil {
				t.Fatal(err)
			}

			err = cache.Store(&storedValue)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		file      string
		expectErr bool
	}{
		{
			name:      "load cache file on temp dir",
			dir:       tmpDir,
			file:      "test1",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := New(tt.dir, tt.file)
			if err != nil {
				t.Fatal(err)
			}

			// remove cache file before test
			if err = cache.Clear(); err != nil {
				t.Fatal(err)
			}

			err = cache.Store(&storedValue)
			if err != nil {
				t.Fatal(err)
			}

			loadedValue := example{}
			err = cache.Load(&loadedValue)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %+v", err)
			}

			if diff := cmp.Diff(storedValue, loadedValue); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}

func TestExpired(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
		want        bool
	}{
		{
			name:        "not expired cache test",
			dir:         tmpDir,
			file:        "test1",
			expiredTime: 3 * time.Minute,
			want:        false,
		},
		{
			name:        "expired cache test",
			dir:         tmpDir,
			file:        "test1",
			expiredTime: 0 * time.Minute,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.dir, tt.file)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Store(&storedValue)
			if err != nil {
				t.Fatal(err)
			}

			if c.Expired(tt.expiredTime) != tt.want {
				t.Errorf("unexpected cache expired or not expired")
			}
		})
	}
}
