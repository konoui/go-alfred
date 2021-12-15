package initialize

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konoui/go-alfred"
	mock "github.com/konoui/go-alfred/update/mock_update"
)

const systemInfoOutput = `{
	"items": [
	  {
		"title": "New version workflow is available!",
		"subtitle": "↩ for update",
		"icon": {
		  "path": "/tmp/go-alfred-data/assets/AlertNoteIcon.icns"
		},
		"autocomplete": "workflow:update",
		"valid": false
	  },
	  {
		"title": "test"
	  }
	]
  }`

const noSystemInfoOutput = `{
	"items": [
	  {
		"title": "test"
	  }
	]
  }`

func TestDisplayUpdateSystemInfo(t *testing.T) {
	tests := []struct {
		name                string
		newVersionAvailable bool
		want                []byte
	}{
		{
			name:                "new version available",
			newVersionAvailable: true,
			want:                []byte(systemInfoOutput),
		},
		{
			name:                "no new version available",
			newVersionAvailable: false,
			want:                []byte(noSystemInfoOutput),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := func(w *alfred.Workflow) error {
				w.Append(alfred.NewItem().Title("test"))
				w.Output()
				return nil
			}

			ctrl := gomock.NewController(t)
			mockSource := mock.NewMockUpdaterSource(ctrl)
			mockSource.EXPECT().IsNewVersionAvailable(gomock.Any()).
				Return(tt.newVersionAvailable, nil)
			defer ctrl.Finish()

			outBuffer := new(bytes.Buffer)
			logBuffer := new(bytes.Buffer)
			w := alfred.NewWorkflow(
				alfred.WithUpdater(mockSource),
				alfred.WithOutWriter(outBuffer),
				alfred.WithLogWriter(logBuffer),
				alfred.WithInitializers(
					NewAutoUpdateChecker(2*time.Second),
				),
			)
			exitCode := w.Run(app)
			if exitCode != 0 {
				t.Errorf("unexpected exit code %d", exitCode)
			}

			got := outBuffer.Bytes()
			if diff := alfred.DiffOutput(tt.want, got); diff != "" {
				t.Errorf("-want/+got %s", diff)
			}
		})
	}
}

func TestAutoUpdateBackgroundJob(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		cmd        string
		wantLogMsg string
	}{
		{
			name:       "do update",
			wantErr:    false,
			cmd:        "echo",
			wantLogMsg: ArgWorkflowUpdate,
		},
		{
			name:       "update command failed due to ls update:workflow",
			wantErr:    true,
			cmd:        "ls",
			wantLogMsg: "background-updater job failed due to",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := make([]string, len(os.Args))
			copy(tmp, os.Args)
			// overwrite os.Args for test
			os.Args = []string{tt.cmd, ArgWorkflowUpdate}
			osExecutable = func() (string, error) {
				return tt.cmd, nil
			}
			osExit = func(code int) {
				if code != 0 {
					t.Fatalf("osExit status code is %d", code)
				}
			}
			defer func() {
				os.Args = tmp
				osExecutable = os.Executable
				osExit = os.Exit
			}()
			newVersionAvailable := true
			ctrl := gomock.NewController(t)
			mockSource := mock.NewMockUpdaterSource(ctrl)
			mockUpdater := mock.NewMockUpdater(ctrl)
			mockSource.EXPECT().IsNewVersionAvailable(gomock.Any()).
				Return(newVersionAvailable, nil).AnyTimes()
			// Note: AnyTimes() is required since IfNewVersionAvailable and Update are called on job process not test process
			mockSource.EXPECT().IfNewVersionAvailable().
				Return(mockUpdater).AnyTimes()
			mockUpdater.EXPECT().Update(gomock.Any()).
				Return(nil).AnyTimes()
			defer ctrl.Finish()

			outBuffer := new(bytes.Buffer)
			logBuffer := new(bytes.Buffer)
			w := alfred.NewWorkflow(
				alfred.WithUpdater(mockSource),
				alfred.WithOutWriter(outBuffer),
				alfred.WithLogWriter(logBuffer),
				alfred.WithInitializers(NewAutoUpdater(2*time.Second)),
			)

			exitCode := w.RunSimple(func() error {
				w.Output()
				return nil
			})

			outStr := outBuffer.String()
			logStr := logBuffer.String()
			if (exitCode != 0) != tt.wantErr {
				t.Errorf("unexpected error out=%s log=%s", outStr, logStr)
			}

			scanner := bufio.NewScanner(logBuffer)
			for scanner.Scan() {
				msg := scanner.Text()
				if strings.Contains(msg, tt.wantLogMsg) {
					return
				}
			}
			t.Errorf("want: %s\ngot log messages: %s\n", tt.wantLogMsg, logStr)
		})
	}
}
