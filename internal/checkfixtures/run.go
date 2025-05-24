package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type RunResult struct {
	Stdout        []byte
	Stderr        []byte
	ExitCode      int
	TimedOut      bool
	OutputLimited bool
	Err           error
}

type limitedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
	mu        sync.Mutex
}

func newLimitedBuffer(limit int) *limitedBuffer {
	return &limitedBuffer{limit: limit}
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	original := len(p)
	if b.limit <= 0 {
		return original, nil
	}
	remaining := b.limit - b.buf.Len()
	if remaining <= 0 {
		b.truncated = true
		return original, nil
	}
	if len(p) > remaining {
		b.buf.Write(p[:remaining])
		b.truncated = true
		return original, nil
	}
	b.buf.Write(p)
	return original, nil
}

func (b *limitedBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]byte(nil), b.buf.Bytes()...)
}

func (b *limitedBuffer) Truncated() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.truncated
}

var prlimitPath = func() string {
	path, err := exec.LookPath("prlimit")
	if err != nil {
		return ""
	}
	return path
}()

func findPython() (string, error) {
	return exec.LookPath("python3")
}

func goToolPath() string {
	tool := filepath.Join(runtime.GOROOT(), "bin", "go")
	if _, err := os.Stat(tool); err != nil {
		return "go"
	}
	return tool
}

func commandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

func runCommand(ctx context.Context, runner commandRunner, input []byte, workDir string, opts Options) RunResult {
	runCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	command := runner.Command
	args := runner.Args
	if prlimitPath != "" && opts.MemoryMB > 0 {
		command = prlimitPath
		args = append([]string{
			"--as=" + formatInt(opts.MemoryMB<<20),
			"--core=0",
			"--cpu=10",
			"--",
			runner.Command,
		}, runner.Args...)
	}

	cmd := exec.CommandContext(runCtx, command, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), runner.Env...)
	cmd.Stdin = bytes.NewReader(input)
	stdout := newLimitedBuffer(16 << 20)
	stderr := newLimitedBuffer(64 << 10)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = 250 * time.Millisecond

	err := cmd.Run()
	result := RunResult{
		Stdout:        stdout.Bytes(),
		Stderr:        stderr.Bytes(),
		OutputLimited: stdout.Truncated(),
		Err:           err,
	}
	if runCtx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}
	return result
}

func formatInt(value int64) string {
	if value == 0 {
		return "0"
	}
	var digits [24]byte
	position := len(digits)
	negative := value < 0
	if negative {
		value = -value
	}
	for value > 0 {
		position--
		digits[position] = byte('0' + value%10)
		value /= 10
	}
	if negative {
		position--
		digits[position] = '-'
	}
	return string(digits[position:])
}
