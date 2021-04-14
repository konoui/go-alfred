package alfred

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/konoui/go-alfred/update/mock_update"
)

func Test_updater_IfNewerVersionAvailable(t *testing.T) {
	tests := []struct {
		name      string
		want      bool
		injectErr error
	}{
		{
			name:      "new version available",
			want:      true,
			injectErr: nil,
		},
		{
			name:      "new version unavailable",
			want:      false,
			injectErr: nil,
		},
		{
			name:      "return false if error occurs",
			want:      false,
			injectErr: errors.New("injected error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mock.NewMockUpdaterSource(ctrl)
			m.EXPECT().NewerVersionAvailable("").Return(tt.want, tt.injectErr)
			updater := &updater{
				source: m,
				wf:     NewWorkflow(),
			}
			defer ctrl.Finish()
			if got := updater.NewerVersionAvailable(); !reflect.DeepEqual(got, tt.want) {
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
					source:         mock.NewMockUpdaterSource(nil),
					currentVersion: "v0.0.1",
				},
			},
			want: &updater{
				source:         mock.NewMockUpdaterSource(nil),
				currentVersion: "v0.0.1",
			},
		},
		{
			wf: &Workflow{
				updater: nil,
			},
			want: &nilUpdater{},
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
		name      string
		args      args
		wantErr   bool
		injectErr error
	}{
		{
			name: "update",
			args: args{
				ctx: context.Background(),
			},
			wantErr:   false,
			injectErr: nil,
		},
		{
			name: "update but return error",
			args: args{
				ctx: context.Background(),
			},
			wantErr:   true,
			injectErr: errors.New("injected error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSource := mock.NewMockUpdaterSource(ctrl)
			mockUpdater := mock.NewMockUpdater(ctrl)
			mockUpdater.EXPECT().Update(tt.args.ctx).Return(tt.injectErr)
			mockSource.EXPECT().NewerVersionAvailable("").Return(true, nil)
			mockSource.EXPECT().IfNewerVersionAvailable("").Return(mockUpdater)
			updater := &updater{
				source: mockSource,
				wf:     NewWorkflow(),
			}
			defer ctrl.Finish()

			if err := updater.Update(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("updater.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
