// Command gate is Rhiza's own local ship gate (rhiza.ship.localGate = make
// gate). From RHIZA_GATE_BASE it runs go test on the Go packages the branch
// changes and their reverse dependencies, cmd/rhiza's testscripts only for the
// changed txtar files and those the change set's leaf acceptance checks name,
// then make self-check. The full suite is left to CI, and runs here only when
// the base is unknown or a file that shapes every test changed.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	dry := flag.Bool("dry-run", false, "print the plan and run nothing")
	flag.Parse()
	if err := run(".", *dry, os.Getenv, os.Stdout, acceptance); err != nil {
		fmt.Fprintln(os.Stderr, "gate:", err)
		os.Exit(1)
	}
}

func run(root string, dry bool, env func(string) string, out io.Writer, accept func(change string) []string) error {
	base, head := env("RHIZA_GATE_BASE"), env("RHIZA_GATE_HEAD")
	if head == "" {
		head = "HEAD"
	}
	var plan Plan
	if base == "" {
		plan = Build("", nil, nil, nil)
	} else {
		changed, err := gitLines(root, "diff", "--name-only", "--no-renames", base, head)
		if err != nil {
			return err
		}
		pkgs, err := Discover(root)
		if err != nil {
			return err
		}
		var acc []string
		if c := env("RHIZA_CHANGE"); c != "" {
			acc = accept(c)
		}
		plan = Build(base, changed, pkgs, acc)
	}
	fmt.Fprintf(out, "gate: %s\n", plan.Reason)
	for _, s := range plan.Steps {
		fmt.Fprintf(out, "  %s\n", s)
	}
	if dry {
		return nil
	}
	for _, s := range plan.Steps {
		fmt.Fprintf(out, "gate: running %s\n", s)
		if err := runStep(root, s, out); err != nil {
			return fmt.Errorf("%s: %w", s, err)
		}
	}
	return nil
}

func runStep(root string, s Step, out io.Writer) error {
	if s.Locked {
		if unlock := lock(out); unlock != nil {
			defer unlock()
		}
	}
	cmd := exec.Command(s.Argv[0], s.Argv[1:]...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, out, out
	return cmd.Run()
}

// lock takes the machine-wide test lock make test uses, so suites queue for
// the CPU-heavy part. Where it cannot be taken nothing is locked.
func lock(out io.Writer) func() {
	file := os.Getenv("TEST_LOCK")
	if file == "" {
		gocache, err := exec.Command("go", "env", "GOCACHE").Output()
		if err != nil {
			return nil
		}
		file = filepath.Join(filepath.Dir(strings.TrimSpace(string(gocache))), "rhiza-test.lock")
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil
	}
	if syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
		fmt.Fprintf(out, "gate: waiting for another make test on this machine (lock: %s)\n", file)
		if syscall.Flock(int(f.Fd()), syscall.LOCK_EX) != nil {
			f.Close()
			return nil
		}
	}
	return func() { syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }
}

func gitLines(root string, args ...string) ([]string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	var lines []string
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}

// acceptance reads the acceptance checks of the change set's leaves from the
// rhiza on $PATH. Failing to read them only narrows which scripts run.
func acceptance(change string) []string {
	b, err := exec.Command("rhiza", "slice", "--change", change, "--json").Output()
	if err != nil {
		return nil
	}
	var s struct {
		Focus []struct {
			Fields map[string]any `json:"fields"`
		} `json:"focus"`
	}
	if json.Unmarshal(b, &s) != nil {
		return nil
	}
	var out []string
	for _, n := range s.Focus {
		if a, ok := n.Fields["acceptance"].(string); ok {
			out = append(out, a)
		}
	}
	return out
}
