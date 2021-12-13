package alfred

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/konoui/go-alfred/update/mock_update"
)

func TestWorkflow_OnInitialize(t *testing.T) {
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
			newerVersionAvailable := true
			ctrl := gomock.NewController(t)
			mockSource := mock.NewMockUpdaterSource(ctrl)
			mockUpdater := mock.NewMockUpdater(ctrl)
			mockSource.EXPECT().NewerVersionAvailable(gomock.Any()).
				Return(newerVersionAvailable, nil).AnyTimes()
			// Note: AnyTimes() is required since IfNewerVersionAvailable and Update are called on job process not test process
			mockSource.EXPECT().IfNewerVersionAvailable().
				Return(mockUpdater).AnyTimes()
			mockUpdater.EXPECT().Update(gomock.Any()).
				Return(nil).AnyTimes()
			defer ctrl.Finish()

			logBuffer := new(bytes.Buffer)
			w := NewWorkflow(
				WithUpdater(mockSource),
				WithOutWriter(io.Discard),
				WithLogWriter(logBuffer),
			)

			if err := w.onInitialize(); (err != nil) != tt.wantErr {
				t.Errorf("Workflow.OnInitialize() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotStr := logBuffer.String()
			scanner := bufio.NewScanner(logBuffer)
			for scanner.Scan() {
				msg := scanner.Text()
				if strings.Contains(msg, tt.wantLogMsg) {
					return
				}
			}
			t.Errorf("want: %s\ngot log messages: %s\n", tt.wantLogMsg, gotStr)
		})
	}
}
