package daemon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	childEnvKey   = "DAEMON_CHILD_FLAG"
	childEnvValue = "child"
)

// ProcessStatus is a type of daemonize results
type ProcessStatus int

const (
	// ParentProcess presents a prarent process for daemonize
	ParentProcess ProcessStatus = iota + 1
	// ChildProcess presents a childor process for daemonsize
	ChildProcess
	// FailedProcess presents daemonize failed
	FailedProcess
)

var (
	// ErrAlreadyRunning presents a process alredy running
	ErrAlreadyRunning = errors.New("the process is already running")
)

// Context presents paremeters of a process
type Context struct {
	ExtraEnv    []string
	PidDir      string
	PidFileName string
}

func (s ProcessStatus) String() string {
	switch s {
	case ParentProcess:
		return "ParentProcess"
	case ChildProcess:
		return "ChildProcess"
	default:
		return "FailedProcess"
	}
}

func (c *Context) Daemonize(cmd *exec.Cmd) (ProcessStatus, error) {
	if c.isParentProcess() {
		if c.IsRunning() {
			return FailedProcess, ErrAlreadyRunning
		}

		// here is a parent process
		if err := c.markAsChildProcess(); err != nil {
			return FailedProcess, err
		}

		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		cmd.SysProcAttr.Setpgid = true
		cmd.Env = append(cmd.Env, c.ExtraEnv...)

		if err := cmd.Start(); err != nil {
			return FailedProcess, err
		}

		// create pid file
		child := cmd.Process
		if err := createPidFile(child.Pid, c.PidFileName, c.PidDir); err != nil {
			_ = child.Kill()
			return FailedProcess, err
		}

		// detach the child process from parent
		if err := child.Release(); err != nil {
			_ = child.Kill()
			return FailedProcess, fmt.Errorf("failed to detach child process: %w", err)
		}
		return ParentProcess, nil
	}

	// here is a child process
	syscall.Umask(0)
	return ChildProcess, nil
}

// IsRunning returns true if the daemon is running
func (c *Context) IsRunning() bool {
	_, err := readPidFile(c.PidFileName, c.PidDir)
	return err == nil
}

// IsChildProcess returns true if the current process is child
func (c *Context) IsChildProcess() bool {
	return !c.isParentProcess()
}

// Terminate kills the child process
func (c *Context) Terminate() error {
	pid, err := readPidFile(c.PidFileName, c.PidDir)
	if err != nil {
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	return proc.Kill()
}

func (c *Context) markAsChildProcess() error {
	keyValue := fmt.Sprintf("%s=%s", childEnvKey, childEnvValue)
	c.ExtraEnv = append(c.ExtraEnv, keyValue)
	return nil
}

func (c *Context) isParentProcess() bool {
	// TODO check c.Env at first
	v := os.Getenv(childEnvKey)
	return v != childEnvValue
}

func createPidFile(pid int, filename, dir string) error {
	tmpfile := filename + ".tmp"
	tmpPidfile := filepath.Join(dir, tmpfile)
	pidfile := filepath.Join(dir, filename)
	data := []byte(fmt.Sprintf("%d", pid))
	err := ioutil.WriteFile(tmpPidfile, data, 0600)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPidfile)

	old, err := readPidFile(filename, dir)
	if err == nil {
		return fmt.Errorf("the process is already running as pid (%d)", old)
	}
	return os.Rename(tmpPidfile, pidfile)
}

func readPidFile(filename, dir string) (int, error) {
	const invalidPid = -1
	pidfile := filepath.Join(dir, filename)
	_, err := os.Stat(pidfile)
	if os.IsNotExist(err) {
		return invalidPid, fmt.Errorf("pid file %s does not exist: %w", pidfile, err)
	}

	v, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return invalidPid, fmt.Errorf("failed to read %s: %w", pidfile, err)
	}

	pid, err := strconv.Atoi(string(v))
	if err != nil {
		_ = os.Remove(pidfile)
		return invalidPid, fmt.Errorf("pid is invalid format: %w", err)
	}

	err = syscall.Kill(pid, syscall.Signal(0))
	if err != nil {
		_ = os.Remove(pidfile)
		return invalidPid, fmt.Errorf("process (%d) does not exist: %w", pid, err)
	}
	return pid, nil
}
