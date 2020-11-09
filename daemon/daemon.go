package daemon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	Name        string
	Args        []string
	PidDir      string
	Dir         string
	Files       []*os.File
	Env         []string
	SysAttr     *syscall.SysProcAttr
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

// Daemonize makes a process as daemon
func (c *Context) Daemonize() (ProcessStatus, error) {
	if c.isParentProcess() {
		if c.IsRunning() {
			return FailedProcess, ErrAlreadyRunning
		}

		// here is a parent process
		err := c.markAsChildProcess()
		if err != nil {
			return FailedProcess, err
		}
		attr := &os.ProcAttr{
			Dir:   c.Dir,
			Env:   c.Env,
			Files: c.Files,
			Sys:   c.SysAttr,
		}
		if attr.Sys == nil {
			attr.Sys = new(syscall.SysProcAttr)
		}
		// overwride sid
		attr.Sys.Setsid = true

		// update command name to abs path as working directory of child may be changed
		asbPath, err := filepath.Abs(c.Name)
		if err != nil {
			return FailedProcess, err
		}
		c.Name = asbPath

		// must set a program name as first argument
		args := append([]string{c.Name}, c.Args...)

		child, err := os.StartProcess(c.Name, args, attr)
		if err != nil {
			return FailedProcess, err
		}

		// create pid file
		err = createPidFile(child.Pid, c.PidFileName, c.PidDir)
		if err != nil {
			_ = child.Kill()
			return FailedProcess, err
		}

		// detach the child process from parent
		if err = child.Release(); err != nil {
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
	c.Env = append(c.Env, keyValue)
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
