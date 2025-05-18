package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type RunOptions struct {
	Timeout     time.Duration
	MaxOutput   int
	MemoryBytes int64
}

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

func runBinary(ctx context.Context, binary string, input []byte, workDir string, env []string, opts RunOptions) RunResult {
	if opts.Timeout <= 0 {
		opts.Timeout = 2 * time.Second
	}
	if opts.MaxOutput <= 0 {
		opts.MaxOutput = 8 << 20
	}

	runCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	var command string
	var args []string
	if prlimitPath != "" && opts.MemoryBytes > 0 {
		command = prlimitPath
		args = append(args,
			"--as="+itoa64(opts.MemoryBytes),
			"--core=0",
			"--cpu=10",
			"--",
			binary,
		)
	} else {
		command = binary
	}

	cmd := exec.CommandContext(runCtx, command, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = bytes.NewReader(input)
	stdout := newLimitedBuffer(opts.MaxOutput)
	stderr := newLimitedBuffer(64 << 10)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = 200 * time.Millisecond

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

func instrumentedEnv(covDir string) []string {
	return []string{"GOCOVERDIR=" + covDir}
}

func goToolPath() string {
	tool := filepath.Join(runtime.GOROOT(), "bin", "go")
	if _, err := os.Stat(tool); err != nil {
		return "go"
	}
	return tool
}

func goTool(args ...string) *exec.Cmd {
	return exec.Command(goToolPath(), args...)
}

func isGoToolAvailable() bool {
	return goTool("version").Run() == nil
}

func itoa64(value int64) string {
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

func normalizeOutput(data []byte) string {
	return strings.Join(strings.Fields(string(data)), " ")
}
