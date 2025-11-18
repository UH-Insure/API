package compiler

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

// Runner wraps all execution logic
type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

//
// --- Generic command runner ---
//

func (r *Runner) runCmd(timeout time.Duration, cmdName string, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = "/work"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", "", err
	}

	err = cmd.Wait()

	// Detect timeout
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), ctx.Err()
	}

	return stdout.String(), stderr.String(), err
}

//
// --- Cryptol ---
//

func (r *Runner) RunCryptol(inputPath string) (string, string, error) {
	// cryptol <file>
	return r.runCmd(
		30*time.Second,
		"cryptol",
		filepath.Base(inputPath),
	)
}

//
// --- SAW ---
//

func (r *Runner) RunSAW(inputPath string) (string, string, error) {
	// saw <file>
	return r.runCmd(
		60*time.Second,
		"saw",
		filepath.Base(inputPath),
	)
}

//
// --- Compile C → LLVM bitcode ---
//

func (r *Runner) CompileCToBitcode(cPath string) (string, string, error) {
	bcOut := filepath.Base(cPath) + ".bc"

	stdout, stderr, err := r.runCmd(
		30*time.Second,
		"clang-16",
		"-emit-llvm",
		"-c",
		filepath.Base(cPath),
		"-o",
		bcOut,
	)

	return stdout, stderr, err
}

//
// --- Multi-dispatch ---
//

func (r *Runner) Execute(tool string, inputPath string) (string, string, error) {
	switch tool {
	case "cryptol":
		return r.RunCryptol(inputPath)

	case "saw":
		return r.RunSAW(inputPath)

	case "c":
		return r.CompileCToBitcode(inputPath)

	default:
		return "", "", exec.ErrNotFound
	}
}