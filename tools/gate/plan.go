package main

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Step is one command of the gate. Locked steps run under the machine-wide
// test lock, as make test runs cmd/rhiza's scripts.
type Step struct {
	Argv   []string
	Locked bool
}

func (s Step) String() string {
	parts := make([]string, len(s.Argv))
	for i, a := range s.Argv {
		if strings.ContainsAny(a, " \t'()|$^*") {
			a = "'" + a + "'"
		}
		parts[i] = a
	}
	line := strings.Join(parts, " ")
	if s.Locked {
		line += "  (under the test lock)"
	}
	return line
}

// Plan is what the gate will run, and why.
type Plan struct {
	Reason string
	Steps  []Step
}

var (
	selfCheck = Step{Argv: []string{"make", "self-check"}}
	fullSteps = []Step{{Argv: []string{"make", "test"}}, selfCheck}

	// A change to one of these can alter how every test runs, so the whole
	// suite runs.
	fullTriggers = map[string]bool{
		"go.mod": true, "go.sum": true, "Makefile": true,
		"cmd/rhiza/main_test.go": true,
	}
	scriptDir  = "cmd/rhiza/testdata/script/"
	cmdRhiza   = "cmd/rhiza"
	scriptName = regexp.MustCompile(`TestScripts/(?:\(([^)]*)\)|([A-Za-z0-9_.-]+))`)
)

// ScriptNames returns the testscript names the acceptance texts run, as in
// "-run TestScripts/foo" or "-run 'TestScripts/(foo|bar)$'".
func ScriptNames(acceptance []string) []string {
	var names []string
	for _, a := range acceptance {
		for _, m := range scriptName.FindAllStringSubmatch(a, -1) {
			if m[1] != "" {
				for _, n := range strings.Split(m[1], "|") {
					names = append(names, strings.Trim(n, "^$ "))
				}
			} else {
				names = append(names, m[2])
			}
		}
	}
	return names
}

// Pkg is one Go package of the module: its directory relative to the repo
// root ("." for the root) and the directories of the packages its code and
// tests import.
type Pkg struct {
	Dir     string
	Imports []string
}

// Build plans the gate for the files changed since the base. An empty base
// means nothing is known about the branch, so the whole suite runs.
func Build(base string, changed []string, pkgs []Pkg, acceptance []string) Plan {
	if base == "" {
		return Plan{"RHIZA_GATE_BASE is empty: no base to scope by", fullSteps}
	}
	var dirs []string
	for _, p := range pkgs {
		dirs = append(dirs, p.Dir)
	}
	changedPkg := map[string]bool{}
	var names []string
	for _, f := range changed {
		if fullTriggers[f] || strings.HasPrefix(f, "tools/gate/") {
			return Plan{f + " changed: it can change every test", fullSteps}
		}
		if strings.HasPrefix(f, scriptDir) && strings.HasSuffix(f, ".txtar") && !strings.Contains(f[len(scriptDir):], "/") {
			names = append(names, strings.TrimSuffix(f[len(scriptDir):], ".txtar"))
		}
		d := owner(f, dirs)
		if d == "" {
			if strings.HasSuffix(f, ".go") {
				return Plan{f + " is Go code outside any package the tool can see", fullSteps}
			}
			if strings.HasPrefix(f, ".rhiza/") {
				continue
			}
			// Tests read files outside their packages (AGENTS.md, agents/,
			// docs/, .github/, .goreleaser.yaml), so which ones a change
			// affects cannot be worked out: the whole suite runs.
			return Plan{f + " is outside every package, and tests may read it", fullSteps}
		}
		changedPkg[d] = true
	}

	// Close over reverse dependencies.
	rev := map[string][]string{}
	for _, p := range pkgs {
		for _, i := range p.Imports {
			rev[i] = append(rev[i], p.Dir)
		}
	}
	affected := map[string]bool{}
	queue := []string{}
	for d := range changedPkg {
		affected[d] = true
		queue = append(queue, d)
	}
	for len(queue) > 0 {
		d := queue[0]
		queue = queue[1:]
		for _, r := range rev[d] {
			if !affected[r] {
				affected[r] = true
				queue = append(queue, r)
			}
		}
	}

	var list []string
	for d := range affected {
		if d != cmdRhiza {
			list = append(list, d)
		}
	}
	sort.Strings(list)
	plan := Plan{Reason: fmt.Sprintf("%d changed file(s) since %s", len(changed), short(base))}
	if len(list) > 0 {
		argv := []string{"go", "test", "-trimpath"}
		for _, d := range list {
			argv = append(argv, rel(d))
		}
		plan.Steps = append(plan.Steps, Step{Argv: argv})
	}
	if affected[cmdRhiza] {
		plan.Steps = append(plan.Steps, Step{Argv: []string{"go", "test", "-trimpath", "./cmd/rhiza", "-skip", "^TestScripts$"}, Locked: true})
		names = append(names, ScriptNames(acceptance)...)
		names = unique(names)
		if len(names) > 0 {
			plan.Steps = append(plan.Steps, Step{
				Argv:   []string{"go", "test", "-trimpath", "./cmd/rhiza", "-run", "TestScripts/(" + strings.Join(names, "|") + ")$"},
				Locked: true,
			})
		}
	}
	plan.Steps = append(plan.Steps, selfCheck)
	return plan
}

func rel(dir string) string {
	if dir == "." {
		return "."
	}
	return "./" + dir
}

func short(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}

func unique(in []string) []string {
	sort.Strings(in)
	var out []string
	for i, s := range in {
		if i == 0 || s != in[i-1] {
			out = append(out, s)
		}
	}
	return out
}

// owner is the deepest package directory holding the file.
func owner(file string, dirs []string) string {
	best := ""
	for d := path.Dir(file); ; d = path.Dir(d) {
		for _, p := range dirs {
			if p == d {
				best = p
				return best
			}
		}
		if d == "." || d == "/" {
			return ""
		}
	}
}

// Discover lists the module's packages with the directories they import
// (code, in-package tests and external tests alike).
func Discover(root string) ([]Pkg, error) {
	run := func(args ...string) (string, error) {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = root
		out, err := cmd.Output()
		return string(out), err
	}
	mod, err := run("go", "list", "-m")
	if err != nil {
		return nil, fmt.Errorf("go list -m: %w", err)
	}
	mod = strings.TrimSpace(mod)
	out, err := run("go", "list", "-f", `{{.ImportPath}}|{{join .Imports " "}} {{join .TestImports " "}} {{join .XTestImports " "}}`, "./...")
	if err != nil {
		return nil, fmt.Errorf("go list ./...: %w", err)
	}
	dir := func(imp string) (string, bool) {
		if imp == mod {
			return ".", true
		}
		if strings.HasPrefix(imp, mod+"/") {
			return strings.TrimPrefix(imp, mod+"/"), true
		}
		return "", false
	}
	var pkgs []Pkg
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		ip, imps, ok := strings.Cut(line, "|")
		if !ok {
			continue
		}
		d, ok := dir(ip)
		if !ok {
			continue
		}
		p := Pkg{Dir: d}
		for _, i := range strings.Fields(imps) {
			if id, ok := dir(strings.TrimSuffix(i, "_test")); ok {
				p.Imports = append(p.Imports, id)
			}
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}
