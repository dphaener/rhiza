package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// TestRunConcurrentDirs checks that two run() calls proceed at once in one
// process, each answering for its own working directory and environment:
// two goroutines drive two scratch repos, with different RHIZA_ACTOR
// values and HOMEs, many times each, and every call's output names its own
// repo's nodes and actor, never the other's. Run it under -race.
func TestRunConcurrentDirs(t *testing.T) {
	type repo struct {
		name, dir string
		env       []string
	}
	t.Chdir(t.TempDir())
	var repos []repo
	for _, name := range []string{"alpha", "beta"} {
		home, dir := t.TempDir(), t.TempDir()
		if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, out)
		}
		// Each HOME's git config names its own person, so git run in the
		// wrong HOME records the wrong one.
		gitconfig := "[user]\n\tname = " + name + "\n\temail = " + name + "@example.com\n"
		if err := os.WriteFile(filepath.Join(home, ".gitconfig"), []byte(gitconfig), 0o644); err != nil {
			t.Fatal(err)
		}
		env := []string{"HOME=" + home, "GIT_CONFIG_NOSYSTEM=1", "PATH=" + os.Getenv("PATH"), "RHIZA_TEST_NOW=2026-10-24T12:00:00Z", "RHIZA_TEST_SEED=1"}
		r := repo{name, dir, env}
		for _, args := range [][]string{
			{"init"},
			{"add", "actor", name + "-actor", "--record"},
			{"add", "capability", name + "-store", "--field", "statement=Keeps " + name + "'s log", "--record"},
		} {
			var stdout, stderr bytes.Buffer
			if code := run(args, dir, "", strings.NewReader(""), &stdout, &stderr, env); code != 0 {
				t.Fatalf("%s: rhiza %s: exit %d\n%s%s", name, strings.Join(args, " "), code, stdout.String(), stderr.String())
			}
		}
		r.env = append(r.env, "RHIZA_ACTOR="+name+"-actor")
		repos = append(repos, r)
	}

	const rounds = 25
	var wg sync.WaitGroup
	for i, r := range repos {
		other := repos[1-i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			rhiza := func(args ...string) (int, string) {
				var stdout, stderr bytes.Buffer
				code := run(args, r.dir, "", strings.NewReader(""), &stdout, &stderr, r.env)
				return code, stdout.String() + stderr.String()
			}
			for n := range rounds {
				ev := fmt.Sprintf("%s-ok-%d", r.name, n)
				if code, out := rhiza("add", "evidence", ev, "--field", "statement="+r.name+" keeps its log", "--field", "result=PASS", "--record"); code != 0 {
					t.Errorf("%s: add evidence %s: exit %d\n%s", r.name, ev, code, out)
					return
				}
				code, out := rhiza("show", ev, "--json")
				if code != 0 || !strings.Contains(out, `"actor": "`+r.name+`-actor"`) || !strings.Contains(out, `"handle": "`+ev+`"`) {
					t.Errorf("%s: show %s: exit %d, want its own actor %s-actor\n%s", r.name, ev, code, r.name, out)
				}
				code, out = rhiza("show", r.name+"-store")
				if code != 0 || !strings.Contains(out, "Keeps "+r.name+"'s log") || strings.Contains(out, other.name) {
					t.Errorf("%s: show %s-store: exit %d\n%s", r.name, r.name, code, out)
				}
				if code, out := rhiza("show", other.name+"-store"); code != 4 {
					t.Errorf("%s: show %s-store found the other repo's node: exit %d\n%s", r.name, other.name, code, out)
				}
			}
		}()
	}
	wg.Wait()
}

// guardedDirs are the packages on run()'s path that read nothing of the
// process: TestNoProcessGlobals parses their non-test Go files.
var guardedDirs = []string{"cmd/rhiza", "internal/api", "internal/run", "internal/daemon", "internal/gosweep", "internal/index"}

// processGlobals are, by import path, the functions and variables that
// read or change the process's working directory, arguments, executable,
// environment or user folders: os's, and path/filepath's Abs, which joins
// a relative path onto os.Getwd (a path a command takes is joined onto
// the Ctx's working directory instead).
var processGlobals = map[string]map[string]bool{
	"os": {
		"Getwd": true, "Chdir": true, "Args": true, "Executable": true,
		"Getenv": true, "LookupEnv": true, "Environ": true, "Setenv": true, "Unsetenv": true,
		"UserCacheDir": true, "UserHomeDir": true, "TempDir": true,
	},
	"path/filepath": {"Abs": true},
}

// processGlobalsAllowed names each remaining site that may read the
// process, by file (slash path from the repository root) and enclosing
// function ("" is the whole file), with its reason.
var processGlobalsAllowed = []struct{ file, fn, reason string }{
	{"cmd/rhiza/main.go", "main", "the process's entry point: it hands run() the process's working directory, arguments, environment and RHIZA_CPUPROFILE"},
	{"cmd/rhiza/ctx.go", "selfPath", "main() alone calls it to find the running binary from os.Args[0] (made absolute with filepath.Abs against the process's own working directory, where os.Args[0] is relative to); run() is given the result"},
	{"internal/daemon/reexec.go", "", "the daemon re-execs its own process image, with the process's own arguments and environment"},
}

// TestNoProcessGlobals checks that no non-test Go file of the guarded
// packages reads the process's working directory (filepath.Abs included),
// arguments, executable, environment or user folders outside
// processGlobalsAllowed; each site it finds is named by file and line. Its
// fixtures show it failing.
func TestNoProcessGlobals(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if found := findProcessGlobals(t, root, guardedDirs); len(found) > 0 {
		t.Errorf("process globals outside the allowlist; take the value from the Ctx or run()'s parameters:\n%s", strings.Join(found, "\n"))
	}

	t.Run("fixtures", func(t *testing.T) {
		fix := t.TempDir()
		write := func(rel, src string) {
			t.Helper()
			p := filepath.Join(fix, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		// An allowlisted file's other functions are not allowlisted, an
		// aliased import is still os, and test files are not parsed.
		write("cmd/rhiza/main.go", "package main\n\nimport \"os\"\n\nfunc main() { _ = os.Getenv(\"A\") }\n\nfunc home() string {\n\treturn os.Getenv(\"HOME\")\n}\n")
		write("internal/index/cache.go", "package index\n\nimport goos \"os\"\n\nvar cacheDir = goos.UserCacheDir\n")
		write("internal/index/cache_test.go", "package index\n\nimport \"os\"\n\nvar _ = os.Getenv(\"X\")\n")
		write("internal/daemon/reexec.go", "package daemon\n\nimport \"os\"\n\nvar _ = os.Args\n")
		// filepath.Abs is flagged under any import name, filepath.Join is
		// not, and selfPath's Abs is allowlisted.
		write("internal/api/apply.go", "package api\n\nimport fp \"path/filepath\"\n\nfunc read(f string) string {\n\tabs, _ := fp.Abs(f)\n\treturn fp.Join(abs, f)\n}\n")
		write("cmd/rhiza/ctx.go", "package main\n\nimport \"path/filepath\"\n\nfunc selfPath() string { p, _ := filepath.Abs(\"x\"); return p }\n\nfunc root(v string) string {\n\tp, _ := filepath.Abs(v)\n\treturn p\n}\n")
		got := findProcessGlobals(t, fix, guardedDirs)
		want := []string{
			"cmd/rhiza/ctx.go:8: filepath.Abs",
			"cmd/rhiza/main.go:8: os.Getenv",
			"internal/api/apply.go:6: filepath.Abs",
			"internal/index/cache.go:5: os.UserCacheDir",
		}
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Errorf("fixtures found\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
	})
}

// findProcessGlobals parses the non-test Go files under dirs (relative to
// root; a missing one is skipped) and returns each use of a processGlobals
// name outside processGlobalsAllowed as "file:line: pkg.Name" (os.Getenv,
// filepath.Abs), sorted.
func findProcessGlobals(t *testing.T, root string, dirs []string) []string {
	t.Helper()
	var found []string
	fset := token.NewFileSet()
	for _, d := range dirs {
		err := filepath.WalkDir(filepath.Join(root, filepath.FromSlash(d)), func(path string, e os.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return filepath.SkipDir
				}
				return err
			}
			if e.IsDir() {
				if e.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
			if err != nil {
				return err
			}
			for _, s := range processGlobalUses(fset, f) {
				if !processGlobalAllowed(rel, s.fn) {
					found = append(found, fmt.Sprintf("%s:%d: %s", rel, s.line, s.name))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	sort.Slice(found, func(i, j int) bool {
		fi, li := splitSite(found[i])
		fj, lj := splitSite(found[j])
		if fi != fj {
			return fi < fj
		}
		return li < lj
	})
	return found
}

// splitSite splits "file:line: pkg.Name" into its file and line.
func splitSite(s string) (string, int) {
	file, rest, _ := strings.Cut(s, ":")
	line, _, _ := strings.Cut(rest, ":")
	n, _ := strconv.Atoi(line)
	return file, n
}

type processGlobalUse struct {
	name, fn string
	line     int
}

// processGlobalUses lists f's selectors of a processGlobals name on its
// import of that name's package, under any import name, called or not,
// each named pkg.Name by the package's last path element, with its line
// and the top-level function it sits in ("" outside any).
func processGlobalUses(fset *token.FileSet, f *ast.File) []processGlobalUse {
	pkgs := map[string]string{} // local import name -> import path
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		if processGlobals[p] == nil {
			continue
		}
		name := path.Base(p)
		if im.Name != nil {
			name = im.Name.Name
		}
		if name != "_" && name != "." {
			pkgs[name] = p
		}
	}
	if len(pkgs) == 0 {
		return nil
	}
	var uses []processGlobalUse
	for _, decl := range f.Decls {
		fn := ""
		if d, ok := decl.(*ast.FuncDecl); ok {
			fn = d.Name.Name
		}
		ast.Inspect(decl, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if x, ok := sel.X.(*ast.Ident); ok {
				if p, ok := pkgs[x.Name]; ok && processGlobals[p][sel.Sel.Name] {
					uses = append(uses, processGlobalUse{path.Base(p) + "." + sel.Sel.Name, fn, fset.Position(sel.Pos()).Line})
				}
			}
			return true
		})
	}
	return uses
}

// processGlobalAllowed reports whether file's function fn is in
// processGlobalsAllowed.
func processGlobalAllowed(file, fn string) bool {
	for _, a := range processGlobalsAllowed {
		if a.file == file && (a.fn == "" || a.fn == fn) {
			return true
		}
	}
	return false
}
