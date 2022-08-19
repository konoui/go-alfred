package update

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func Test_compareVersions(t *testing.T) {
	type args struct {
		v2Str string
		v1Str string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "v1 is newer than v2",
			args: args{
				v1Str: "v1.0.0",
				v2Str: "v0.0.1",
			},
			want: false,
		},
		{
			name: "v2 is newer than v1",
			args: args{
				v1Str: "v0.0.1",
				v2Str: "v1.0.0",
			},
			want: true,
		},
		{
			name: "v1 is invalid format",
			args: args{
				v1Str: "invalid",
				v2Str: "v1.0.0",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "v2 is invalid format",
			args: args{
				v1Str: "v0.0.1",
				v2Str: "invalid",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compareVersions(tt.args.v2Str, tt.args.v1Str)
			if (err != nil) != tt.wantErr {
				t.Errorf("compareVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("compareVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupServer(t *testing.T, path string) *httptest.Server {
	t.Helper()
	h := http.NewServeMux()
	h.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	return httptest.NewServer(h)
}

func Test_donwloadContext(t *testing.T) {
	type args struct {
		ctx  context.Context
		url  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid path",
			args: args{
				ctx:  context.TODO(),
				path: tmpDir,
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			args: args{
				ctx:  context.TODO(),
				url:  "aaaa",
				path: filepath.Join(tmpDir, "test.txt"),
			},
			wantErr: true,
		},
		{
			name: "download",
			args: args{
				ctx:  context.TODO(),
				url:  "",
				path: filepath.Join(tmpDir, "test.txt"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "test.txt"
			ts := setupServer(t, path)
			defer ts.Close()

			if tt.args.url == "" {
				tt.args.url = ts.URL + "/" + path
			}

			if err := donwloadContext(tt.args.ctx, tt.args.url, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("donwloadContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateContext(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid url",
			args: args{
				ctx: context.TODO(),
				url: "aaaa",
			},
			wantErr: true,
		},
		{
			name: "download",
			args: args{
				ctx: context.TODO(),
				url: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// switch command
			openCmd = "ls"
			path := "test.txt"
			ts := setupServer(t, path)
			defer ts.Close()

			if tt.args.url == "" {
				tt.args.url = ts.URL + "/" + path
			}

			if err := updateContext(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("updateContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
