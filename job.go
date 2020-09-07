package alfred

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/konoui/go-alfred/daemon"
)

const jobDirKey = "job-dir"
const (
	// ForegroundParentJob presents running process in forground that means parent process
	ForegroundParentJob JobProcess = JobProcess(daemon.ParentProcess)
	// BackgroundChildJob presents running process as job in background that means child process
	BackgroundChildJob JobProcess = JobProcess(daemon.ChildProcess)
	// FailedJob presents failed to start Job
	FailedJob JobProcess = JobProcess(daemon.FailedProcess)
)

// JobProcess is a type of job
type JobProcess int

func (j JobProcess) String() string {
	switch j {
	case ForegroundParentJob:
		return "ForegroundParentJob"
	case BackgroundChildJob:
		return "BackgroundChildJob"
	default:
		return "FailedJob"
	}
}

// Job is context
type Job struct {
	name      string
	dir       string
	daemonCtx *daemon.Context
	wf        *Workflow
	logging   bool
}

// Job creates new job. name parameter means pid file
func (w *Workflow) Job(name string) *Job {
	c := new(daemon.Context)
	c.PidFileName = name + ".pid"
	return &Job{
		name:      name,
		daemonCtx: c,
		wf:        w,
		dir:       w.getJobDir(),
	}
}

func (w *Workflow) getJobDir() string {
	dir, ok := w.dirs[jobDirKey]
	if ok {
		return dir
	}
	return "./"
}

// SetJobDir set job data directory
func (w *Workflow) SetJobDir(dir string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		return
	}
	w.dirs[jobDirKey] = dir
	return
}

// ListJobs return jobs managed by workflow
func (w *Workflow) ListJobs() []*Job {
	const ext = ".pid"
	dir := w.getJobDir()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		w.logger.Printf("Invalid directory %s\n", dir)
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
		w.logger.Printf("Found a job %v\n", job)
		jobs = append(jobs, job)
	}
	return jobs
}

// Name returns job name
func (j *Job) Name() string {
	return j.name
}

// Logging redirects output of stdout/err to log file
func (j *Job) Logging() *Job {
	j.logging = true
	return j
}

// Start behaves as fork if run a self program. If run a external command, it behaves as fork/exec
func (j *Job) Start(cmdName string, args ...string) (JobProcess, error) {
	absPath, err := exec.LookPath(cmdName)
	if err != nil {
		return FailedJob, err
	}
	j.daemonCtx.Name = absPath
	j.daemonCtx.Args = args
	j.daemonCtx.Dir = j.dir
	j.daemonCtx.Files = j.files()
	j.daemonCtx.Env = os.Environ()

	ret, err := j.daemonCtx.Daemonize()
	if err != nil {
		return FailedJob, err
	}
	return JobProcess(ret), nil
}

func (j *Job) files() []*os.File {
	if !j.logging {
		return nil
	}

	p := filepath.Join(j.dir, j.name+".log")
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return []*os.File{
		nil,
		f,
		f,
	}
}

// StartWithExit outputs items and continues only a child process.
// It is helpful to run next instructions as daemon
func (j *Job) StartWithExit(cmdName string, args ...string) *Workflow {
	ret, err := j.Start(cmdName, args...)
	if err != nil {
		j.wf.Fatal("Failed to start a job", err.Error())
	}
	if ret == ForegroundParentJob {
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
