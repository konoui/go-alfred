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
		name          string
		want          bool
		setupMockFunc func(source *mock.MockUpdaterSource)
	}{
		{
			name: "new version available",
			want: true,
			setupMockFunc: func(source *mock.MockUpdaterSource) {
				source.EXPECT().NewerVersionAvailable(gomock.Any()).Return(true, nil)
			},
		},
		{
			name: "new version unavailable",
			want: false,
			setupMockFunc: func(source *mock.MockUpdaterSource) {
				source.EXPECT().NewerVersionAvailable(gomock.Any()).Return(false, nil)
			},
		},
		{
			name: "return false if error occurs",
			setupMockFunc: func(source *mock.MockUpdaterSource) {
				source.EXPECT().NewerVersionAvailable(gomock.Any()).Return(false, errors.New("injected error"))
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mock.NewMockUpdaterSource(ctrl)
			tt.setupMockFunc(m)
			updater := &updater{
				source: m,
				wf:     NewWorkflow(),
			}
			defer ctrl.Finish()
			if got := updater.NewerVersionAvailable(context.TODO()); !reflect.DeepEqual(got, tt.want) {
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
	tests := []struct {
		name          string
		setupMockFunc func(source *mock.MockUpdaterSource, updater *mock.MockUpdater)
		wantErr       bool
	}{
		{
			name: "update",
			setupMockFunc: func(source *mock.MockUpdaterSource, updater *mock.MockUpdater) {
				updater.EXPECT().
					Update(gomock.Any()).
					Return(nil)
				source.EXPECT().NewerVersionAvailable(gomock.Any()).Return(true, nil)
				source.EXPECT().IfNewerVersionAvailable().Return(updater)
			},
			wantErr: false,
		},
		{
			name: "update but return error",
			setupMockFunc: func(source *mock.MockUpdaterSource, updater *mock.MockUpdater) {
				updater.EXPECT().
					Update(gomock.Any()).
					Return(errors.New("injected error"))
				source.EXPECT().NewerVersionAvailable(gomock.Any()).Return(true, nil)
				source.EXPECT().IfNewerVersionAvailable().Return(updater)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSource := mock.NewMockUpdaterSource(ctrl)
			mockUpdater := mock.NewMockUpdater(ctrl)
			tt.setupMockFunc(mockSource, mockUpdater)
			updater := &updater{
				source: mockSource,
				wf:     NewWorkflow(),
			}
			defer ctrl.Finish()

			if err := updater.Update(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("updater.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
