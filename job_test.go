package alfred

// func TestJob_Start(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		cmd  *exec.Cmd
// 		want JobProcess
// 	}{
// 		{
// 			name: "ls",
// 			cmd:  exec.Command("ls", "/bin/ls"),
// 			want: JobStarter,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wf := testWorkflow()
// 			j := wf.Job("test")

// 			got, err := j.Logging().Start(tt.cmd)
// 			if err != nil {
// 				t.Fatalf("Job.Start() error=%v", err)
// 			}

// 			if got != tt.want {
// 				t.Errorf("Job.Start() = %v, want %v", got, tt.want)
// 			}

// 			if err := tt.cmd.Wait(); err != nil {
// 				t.Fatalf("command error: %v", err)
// 			}

// 			gotLog := getLog(t)
// 			if !strings.Contains(gotLog, "/bin/ls") {
// 				t.Errorf("unexpected out %v", got)
// 			}
// 		})
// 	}
// }

// func getLog(t *testing.T) string {
// 	cmd := exec.Command("/usr/bin/log", "show", "--style", "compact", "--info", "--predicate", "'process == \"logger\"'", "--last", "1m")
// 	buf := &bytes.Buffer{}
// 	cmd.Stdout = buf
// 	cmd.Env = os.Environ()
// 	if err := cmd.Run(); err != nil {
// 		t.Fatalf("getLog failed %v", err)
// 	}
// 	return buf.String()
// }
