package alfred

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/konoui/go-alfred/internal/asl"
	"github.com/konoui/go-alfred/internal/daemon"
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

const pidExt = ".pid"

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

type Jobber interface {
	Name() string
	Logging() *Job
	Start(cmd *exec.Cmd) (JobProcess, error)
	IsJob() bool
	IsRunning() bool
	Terminate() error
}

// Job is context
type Job struct {
	name      string
	daemonCtx *daemon.Context
	wf        *Workflow
	logging   bool
}

func (w *Workflow) getJobDir() string {
	dir, err := GetWorkflowDir()
	if err != nil {
		w.sLogger().Warnf("using tmp dir %s for job dir as %s", tmpDir, err)
		return tmpDir
	}

	jobDir := filepath.Join(dir, "jobs")
	if !PathExists(jobDir) {
		if err := os.MkdirAll(jobDir, os.ModePerm); err != nil {
			w.sLogger().Warnf("cannot create job dir due to %s", err)
			return tmpDir
		}
	}
	return jobDir
}

// Job creates new job. name parameter means pid file
func (w *Workflow) Job(name string) *Job {
	c := new(daemon.Context)
	c.PidFileName = name + pidExt
	c.PidDir = w.getJobDir()
	return &Job{
		name:      name,
		daemonCtx: c,
		wf:        w,
	}
}

// ListJobs returns jobs managed by the workflow
func (w *Workflow) ListJobs() []*Job {
	dir := w.getJobDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		w.sLogger().Infof("invalid directory %s\n", dir)
		return nil
	}

	var jobs []*Job
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		if !strings.HasSuffix(filename, pidExt) {
			continue
		}

		jobName := filename[:len(filename)-len(pidExt)]
		job := w.Job(jobName)
		if !job.IsRunning() {
			continue
		}

		// valid process
		pidfile := filepath.Join(job.daemonCtx.PidDir, job.daemonCtx.PidFileName)
		w.sLogger().Debugf("found a job in %s", pidfile)
		jobs = append(jobs, job)
	}

	w.sLogger().Infof("results: found %d jobs", len(jobs))
	return jobs
}

// Name returns job name
func (j *Job) Name() string {
	return j.name
}

// Loggin enables asl(syslog) logging for the job.
// stdout/stderr of the job process will redirect to asl
func (j *Job) Logging() *Job {
	j.logging = true
	return j
}

// Start behaves as fork if run a self program. If run a external command, it behaves as fork/exec
func (j *Job) Start(cmd *exec.Cmd) (JobProcess, error) {
	mergeEnv(cmd, os.Environ())
	if j.logging {
		l := asl.New()
		cmd.Stdout, cmd.Stderr = l, l
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
		j.wf.sLogger().Errorln("Failed to start a job", err.Error())
		j.wf.Fatal("Failed to start a job", err.Error())
	}
	if ret == JobStarter {
		j.wf.Output()
		osExit(0)
	}
	return j.wf
}

// IsJob returns true if a caller is a job but the caller may no be the job
func (j *Job) IsJob() bool {
	return j.daemonCtx.IsChildProcess()
}

// IsRunning returns true if the job(process) is running
// If caller is a job, it will return false.
func (j *Job) IsRunning() bool {
	// Note ignore case that caller is the job
	return !j.IsJob() && j.daemonCtx.IsRunning()
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
