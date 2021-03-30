package alfred

import (
	"context"
	"reflect"
	"testing"

	"github.com/konoui/go-alfred/update"
)

type MockUpdater struct{}

func newMockUpdater() update.UpdaterSource {
	return &MockUpdater{}
}

func (m *MockUpdater) NewerVersionAvailable(currentVersion string) (bool, error) {
	return true, nil
}

func (m *MockUpdater) IfNewerVersionAvailable(currentVersion string) update.Updater {
	return m
}

func (m *MockUpdater) Update(ctx context.Context) error {
	return nil
}

func Test_updater_IfNewerVersionAvailable(t *testing.T) {
	type args struct {
		currentVersion string
	}
	tests := []struct {
		name    string
		updater updater
		args    args
		want    bool
	}{
		{
			name: "updater",
			updater: updater{
				source: newMockUpdater(),
				wf:     NewWorkflow(),
			},
			args: args{
				currentVersion: "v0.0.1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.updater.NewerVersionAvailable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updater.IfNewerVersionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflow_Updater(t *testing.T) {

	tests := []struct {
		name string
		wf   *Workflow
		want Updater
	}{
		{
			wf: &Workflow{
				updater: &updater{
					source:         newMockUpdater(),
					currentVersion: "v0.0.1",
				},
			},
			want: &updater{
				source:         newMockUpdater(),
				currentVersion: "v0.0.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wf.Updater(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Workflow.Updater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updater_Update(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		updater Updater
		args    args
		wantErr bool
	}{
		{
			name: "update",
			updater: &updater{
				source: newMockUpdater(),
				wf:     NewWorkflow(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.updater.Update(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("updater.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
