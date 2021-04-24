package alfred

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/konoui/go-alfred/daemon"
)

const (
	// obStarter running process in forground that means parent process
	JobStarter JobProcess = JobProcess(daemon.ParentProcess)
	// JJobWorker presents running process as job in background that means child process
	JobWorker JobProcess = JobProcess(daemon.ChildProcess)
	// JobFailed presents failed to start Job
	JobFailed JobProcess = JobProcess(daemon.FailedProcess)
)

// JobProcess is a type of job
type JobProcess int

func (j JobProcess) String() string {
	switch j {
	case JobStarter:
		return "JobStarter"
	case JobWorker:
		return "JobWorker"
	default:
		return "JobFailed"
	}
}

// Job is context
type Job struct {
	name      string
	daemonCtx *daemon.Context
	wf        *Workflow
	logging   bool
}

var tmpDir = os.TempDir()

func (w *Workflow) getJobDir() string {
	dir, err := w.GetWorkflowDir()
	if err != nil {
		w.Logger().Warnln("using tmp dir for job dir as", err)
		return tmpDir
	}

	jobDir := filepath.Join(dir, "jobs")
	if err := os.MkdirAll(jobDir, os.ModePerm); err != nil {
		w.Logger().Warnln("cannot create job dir due to", err)
		return tmpDir
	}
	return jobDir
}

// Job creates new job. name parameter means pid file
func (w *Workflow) Job(name string) *Job {
	c := new(daemon.Context)
	c.PidFileName = name + ".pid"
	c.PidDir = w.getJobDir()
	return &Job{
		name:      name,
		daemonCtx: c,
		wf:        w,
	}
}

// ListJobs return jobs managed by workflow
func (w *Workflow) ListJobs() []*Job {
	const ext = ".pid"
	dir := w.getJobDir()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		w.Logger().Infof("invalid directory %s\n", dir)
		return nil
	}

	var jobs []*Job
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		if !strings.HasSuffix(filename, ext) {
			continue
		}

		jobName := filename[:len(filename)-len(ext)]
		job := w.Job(jobName)
		if !job.IsRunning() {
			continue
		}
		// valid process
		w.Logger().Infof("found a job %v\n", job)
		jobs = append(jobs, job)
	}
	return jobs
}

// Name returns job name
func (j *Job) Name() string {
	return j.name
}

func (j *Job) Logging() *Job {
	j.logging = true
	return j
}

// Start behaves as fork if run a self program. If run a external command, it behaves as fork/exec
func (j *Job) Start(cmd *exec.Cmd) (JobProcess, error) {
	mergeEnv(cmd, os.Environ())
	if j.logging {
		cmd.Stdout, cmd.Stderr = j.files()
	}
	ret, err := j.daemonCtx.Daemonize(cmd)
	if err != nil {
		return JobWorker, err
	}
	return JobProcess(ret), nil
}

// StartWithExit outputs items and continues only a child process.
// It is helpful to run next instructions as daemon
func (j *Job) StartWithExit(cmd *exec.Cmd) *Workflow {
	ret, err := j.Start(cmd)
	if err != nil {
		j.wf.Fatal("Failed to start a job", err.Error())
	}
	if ret == JobStarter {
		j.wf.Output()
		os.Exit(0)
	}
	return j.wf
}

// IsJob returns true if I am a job but I may no be the job
func (j *Job) IsJob() bool {
	return j.daemonCtx.IsChildProcess()
}

// IsRunning returns true if the job(process) is running
func (j *Job) IsRunning() bool {
	return j.daemonCtx.IsRunning()
}

// Terminate kills the job
func (j *Job) Terminate() error {
	return j.daemonCtx.Terminate()
}

// TODO restrict key and value
func mergeEnv(cmd *exec.Cmd, envs []string) {
	contains := func(v string, list []string) bool {
		for _, c := range list {
			if c == v {
				return true
			}
		}
		return false
	}

	for _, env := range envs {
		if contains(env, cmd.Env) {
			continue
		}
		cmd.Env = append(cmd.Env, env)
	}
}

// Logging redirects output of stdout/err to log file
func (j *Job) files() (out, stderr *os.File) {
	absPath, err := filepath.Abs(filepath.Join(j.daemonCtx.PidDir, j.name+".log"))
	if err != nil {
		return nil, nil
	}

	f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	j.wf.Logger().Debugln("job logs will be stored at", absPath)
	if err != nil {
		return nil, nil
	}
	return f, f
}
