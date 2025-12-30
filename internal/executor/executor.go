package executor

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"syscall"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	CpuTime  int64
	Memory   int64
}

func (e *Executor) Execute(ctx context.Context, code string) (*Result, error) {
	var stdout, stderr bytes.Buffer

	args := []string{
		"-Mo",
		"-q",
		"--chroot", "/",
		"--user", "2000",
		"--group", "2000",

		"--rlimit_as", "128",
		"--rlimit_nproc", "16",
		"--rlimit_fsize", "10",
		"--rlimit_cpu", "5",

		"--tmpfsmount", "/app",
		"--cwd", "/app",
		"--bindmount_ro", "/usr:/usr",
		"--bindmount_ro", "/lib:/lib",
		"--bindmount_ro", "/lib64:/lib64",
		"--bindmount_ro", "/bin:/bin",
		"--proc_rw",

		"--time_limit", "7",

		"--", "/usr/bin/python3", "-c", code,
	}

	cmd := exec.CommandContext(ctx, "nsjail", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}

	var cpuMs int64
	var maxRss int64
	if cmd.ProcessState != nil {
		if rusage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage); ok {
			uMs := rusage.Utime.Sec*1000 + int64(rusage.Utime.Usec/1000)
			sMs := rusage.Stime.Sec*1000 + int64(rusage.Stime.Usec/1000)
			cpuMs = uMs + sMs
			maxRss = rusage.Maxrss
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		CpuTime:  cpuMs,
		Memory:   maxRss,
	}, nil
}
