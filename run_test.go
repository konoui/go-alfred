package alfred

import (
	"fmt"
	"io"
	"testing"
)

type testInitializer struct{}

func (*testInitializer) Condition() bool { return true }
func (*testInitializer) Initialize(w *Workflow) error {
	return fmt.Errorf("initializer error")
}

func TestWorkflow_Run(t *testing.T) {
	type args struct {
		fn          func(*Workflow) error
		initializer []Initializer
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "run and exit 0",
			args: args{
				fn: func(w *Workflow) error {
					w.Output()
					return nil
				},
			},
			want: 0,
		},
		{
			name: "run but error occurs",
			args: args{
				fn: func(w *Workflow) error {
					return fmt.Errorf("error occurs on fn")
				},
			},
			want: 1,
		},
		{
			name: "run but panic on fn",
			args: args{
				fn: func(w *Workflow) error {
					panic("panic on fn")
				},
			},
			want: 1,
		},
		{
			name: "initialize returns error",
			args: args{
				fn:          func(w *Workflow) error { return nil },
				initializer: []Initializer{new(testInitializer)},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorkflow(
				WithLogWriter(io.Discard),
			)
			got := w.Run(tt.args.fn, tt.args.initializer...)
			if tt.want != got {
				t.Errorf("want: %d got: %d", tt.want, got)
			}
		})
	}
}

func TestWorkflow_RunSimple(t *testing.T) {
	wantExitCode := 0
	w := NewWorkflow()
	fn := func() error {
		w.Output()
		return nil
	}
	got := w.RunSimple(fn)
	if wantExitCode != got {
		t.Errorf("want: %d got: %d", wantExitCode, got)
	}
}
