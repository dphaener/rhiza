//go:build unix

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// detachProgram is api.Ctx.Detach: it starts p in a session of its own, so
// it outlives this process and its terminal, with no input and its output
// appended to logPath. In a test run RHIZA_TEST_DETACH_LOG, when set, names
// a file that receives the environment set and the command line instead and nothing starts.
func detachProgram(p apiProgram, logPath string) error {
	if rec := os.Getenv("RHIZA_TEST_DETACH_LOG"); rec != "" && os.Getenv("RHIZA_TEST_NOW") != "" {
		f, err := os.OpenFile(rec, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(strings.Join(append(append([]string{}, p.Env...), append([]string{"rhiza"}, p.Args...)...), " ") + "\n")
		return err
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return err
	}
	log, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer log.Close()
	cmd := exec.Command(p.Name, p.Args...)
	cmd.Dir = p.Dir
	cmd.Env = programEnv(os.Environ(), p.Env, p.Unset)
	cmd.Stdout, cmd.Stderr = log, log
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
