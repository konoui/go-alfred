package alfred

import (
	"os"
	"testing"

	"golang.org/x/text/unicode/norm"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "convert",
			query: norm.NFD.String("ブックマーク"),
			want:  "ブックマーク",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == tt.query {
				t.Fatal("invalid test")
			}
			if got := Normalize(tt.query); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "true",
			value: "true",
			want:  true,
		},
		{
			name:  "enable",
			value: "enable",
			want:  true,
		},
		{
			name:  "disable",
			value: "disable",
			want:  false,
		},
		{
			name:  "1",
			value: "1",
			want:  true,
		},
		{
			name:  "0",
			value: "0",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.value); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAutoUpdateWorkflowEnabled(t *testing.T) {
	tests := []struct {
		name  string
		want  bool
		value string
	}{
		{
			name:  "enable explicitly with env",
			value: "true",
			want:  true,
		},
		{
			name:  "enable implicitly without env",
			value: "",
			want:  true,
		},
		{
			name:  "disable explicitly with env",
			value: "false",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv(EnvAutoUpdateWorkflow, tt.value); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Unsetenv(EnvAutoUpdateWorkflow); err != nil {
					t.Fatal(err)
				}
			}()
			if got := IsAutoUpdateWorkflowEnabled(); got != tt.want {
				t.Errorf("IsAutoUpdateWorkflowEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflow_GetWorkflowDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "return err since directory does not exist",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow()
			if err := w.OnInitialize(); err != nil {
				t.Fatal(err)
			}

			got, err := w.GetWorkflowDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("Workflow.GetWorkflowDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Workflow.GetWorkflowDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
