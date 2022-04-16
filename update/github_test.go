package update

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	mock "github.com/konoui/go-alfred/update/mock_github"
)

func repositoryReleaseData(t *testing.T, filepath string) *github.RepositoryRelease {
	t.Helper()
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	repositoryRelease := &github.RepositoryRelease{}
	if err := json.Unmarshal(data, repositoryRelease); err != nil {
		t.Fatal(err)
	}
	return repositoryRelease
}

func Test_getLatestAssetInfo(t *testing.T) {
	type args struct {
		ctx       context.Context
		filepath  string
		injectErr error
	}
	tests := []struct {
		name    string
		args    args
		wantTag string
		wantURL string
		wantErr bool
	}{
		{
			name: "workflow artifact",
			args: args{
				ctx:      context.TODO(),
				filepath: "testdata/workflow-latest-release.json",
			},
			wantTag: "v1.0.0",
			wantURL: "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.alfredworkflow",
			wantErr: false,
		},
		{
			name: "not workflow artifact",
			args: args{
				ctx:      context.TODO(),
				filepath: "testdata/non-workflow-latest-release.json",
			},
			wantTag: "",
			wantURL: "",
			wantErr: true,
		},
		{
			name: "github api error",
			args: args{
				filepath:  "testdata/workflow-latest-release.json",
				injectErr: errors.New("get-latest-release-error"),
			},
			wantTag: "",
			wantURL: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			releaseData := repositoryReleaseData(t, tt.args.filepath)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock.NewMockRepositoriesService(ctrl)
			m.EXPECT().GetLatestRelease(tt.args.ctx, "", "").Return(releaseData, nil, tt.args.injectErr)

			gotTag, gotURL, err := getLatestAssetInfo(tt.args.ctx, m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestAssetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTag != tt.wantTag {
				t.Errorf("getLatestAssetInfo() got = %v, want %v", gotTag, tt.wantTag)
			}
			if gotURL != tt.wantURL {
				t.Errorf("getLatestAssetInfo() got1 = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_gitHubUpdater_newVersionAvailableContext(t *testing.T) {
	type fields struct {
		checkInterval time.Duration
	}
	type args struct {
		ctx            context.Context
		currentVersion string
		filepath       string
		injectErr      error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantURL string
		wantOK  bool
		wantErr bool
	}{
		{
			name: "new version available",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v0.0.1",
			},
			wantOK:  true,
			wantURL: "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.alfredworkflow",
			wantErr: false,
		},
		{
			name: "no new version",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v1000.0.0",
			},
			wantOK:  false,
			wantURL: "",
			wantErr: false,
		},
		{
			name: "cache is not expired",
			fields: fields{
				checkInterval: 100 * time.Hour,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v0.0.1",
			},
			wantOK:  false,
			wantURL: "",
			wantErr: false,
		},
		{
			name: "invalid version format",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "invalid",
			},
			wantOK:  false,
			wantURL: "",
			wantErr: true,
		},
		{
			name: "github api error",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v0.0.1",
				injectErr:      errors.New("get-latest-release-error"),
			},
			wantOK:  false,
			wantURL: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			releaseData := repositoryReleaseData(t, tt.args.filepath)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock.NewMockRepositoriesService(ctrl)
			m.EXPECT().GetLatestRelease(tt.args.ctx, "", "").AnyTimes().Return(releaseData, nil, tt.args.injectErr)
			g := &gitHubUpdater{
				client:         m,
				owner:          "",
				repo:           "",
				currentVersion: tt.args.currentVersion,
				checkInterval:  tt.fields.checkInterval,
			}
			gotOK, gotURL, err := g.newVersionAvailableContext(tt.args.ctx, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("gitHubUpdater.newVersionAvailableContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOK != tt.wantOK {
				t.Errorf("gitHubUpdater.newVersionAvailableContext() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if gotURL != tt.wantURL {
				t.Errorf("gitHubUpdater.newVersionAvailableContext() gotURL = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_gitHubUpdater_NewVersionAvailable(t *testing.T) {
	type fields struct {
		checkInterval time.Duration
	}
	type args struct {
		currentVersion string
		filepath       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "no new version available",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v1000.0.0",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			releaseData := repositoryReleaseData(t, tt.args.filepath)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock.NewMockRepositoriesService(ctrl)
			m.EXPECT().GetLatestRelease(context.TODO(), "", "").AnyTimes().Return(releaseData, nil, nil)
			g := &gitHubUpdater{
				currentVersion: tt.args.currentVersion,
				client:         m,
				checkInterval:  tt.fields.checkInterval,
			}
			got, err := g.IsNewVersionAvailable(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("gitHubUpdater.NewVersionAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gitHubUpdater.NewVersionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gitHubUpdater_Update(t *testing.T) {
	type fields struct {
		checkInterval time.Duration
	}
	type args struct {
		ctx            context.Context
		currentVersion string
		filepath       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no new version available",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "v1000.0.0",
			},
			wantErr: false,
		},
		{
			name: "invalid format error",
			fields: fields{
				checkInterval: 0,
			},
			args: args{
				filepath:       "testdata/workflow-latest-release.json",
				currentVersion: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			releaseData := repositoryReleaseData(t, tt.args.filepath)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock.NewMockRepositoriesService(ctrl)
			m.EXPECT().GetLatestRelease(tt.args.ctx, "", "").AnyTimes().Return(releaseData, nil, nil)
			g := &gitHubUpdater{
				client:         m,
				checkInterval:  tt.fields.checkInterval,
				currentVersion: tt.args.currentVersion,
			}
			err := g.IfNewVersionAvailable().Update(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("gitHubUpdater.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGitHubSource(t *testing.T) {
	type args struct {
		owner          string
		repo           string
		currentVersion string
		opts           []Option
	}
	tests := []struct {
		name string
		args args
		want UpdaterSource
	}{
		{
			name: "case 1",
			args: args{
				owner:          "owner",
				repo:           "repo",
				currentVersion: "v0.0.1",
				opts: []Option{
					WithCheckInterval(0),
				},
			},
			want: &gitHubUpdater{
				client:         github.NewClient(nil).Repositories,
				owner:          "owner",
				repo:           "repo",
				currentVersion: "v0.0.1",
				checkInterval:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGitHubSource(tt.args.owner, tt.args.repo, tt.args.currentVersion, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitHubSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
