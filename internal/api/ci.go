package api

// This file holds the ci.install operation: it writes a GitHub workflow that
// runs rhiza check on pull requests and pushes to main, and marks the
// rendered folder as generated in .gitattributes so reviews fold it away.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// The files ci.install writes, relative to the project root, with slashes.
const (
	workflowPath      = ".github/workflows/rhiza.yml"
	gitattributesPath = ".gitattributes"
)

// gitattribute is the .gitattributes line that marks rendered files as
// generated.
const gitattribute = renderedDir + "/** linguist-generated=true"

// workflowTemplate is the workflow; %s is the version go install fetches.
// The runner and action versions match the Rhiza repository's own CI.
const workflowTemplate = `# Written by ` + "`rhiza ci install`" + `. Rerun ` + "`rhiza ci install --force`" + ` to
# replace it with the current template.
#
# It checks the Rhiza model on every pull request against its base branch
# and on every push to main against the commit before the push. The go
# install step needs read access to github.com/dphaener/rhiza: while that
# repository is private, set GOPRIVATE and give git a token that can read it.
name: Rhiza

on:
  push:
    branches: [ main ]
  pull_request:

# A newer push to a pull request cancels its in-progress run. A push to main
# never cancels: its check is incremental from github.event.before, so a
# cancelled run would leave that push's changes unchecked.
concurrency:
  group: rhiza-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: ${{ github.event_name == 'pull_request' }}

jobs:
  check:
    runs-on: ubuntu-26.04

    steps:
    - uses: actions/checkout@v7
      with:
        fetch-depth: 0

    - name: Set up Go
      uses: actions/setup-go@v7
      with:
        go-version: stable

    - name: Install rhiza
      run: go install github.com/dphaener/rhiza/cmd/rhiza@%s

    - name: Check against the base branch
      if: github.event_name == 'pull_request'
      run: rhiza check --base origin/${{ github.base_ref }}

    - name: Check against the commit before the push
      if: github.event_name == 'push' && github.event.before != '0000000000000000000000000000000000000000'
      run: rhiza check --base ${{ github.event.before }}

    # A branch's first push has no commit before it: github.event.before is
    # all zeros, so check runs with its default base.
    - name: Check a first push
      if: github.event_name == 'push' && github.event.before == '0000000000000000000000000000000000000000'
      run: rhiza check
`

// CIInstallRequest asks for the workflow, overwriting a different one only
// with Force.
type CIInstallRequest struct {
	Force bool `json:"force,omitempty" desc:"Overwrite an existing workflow that differs from the template"`
}

// CIInstallResponse says which files ci.install wrote and which were
// already current. Paths are relative to the project root, with slashes,
// sorted.
type CIInstallResponse struct {
	Written   []string `json:"written" desc:"The files written: new, or different from what was on disk"`
	Unchanged []string `json:"unchanged" desc:"The files that were already current"`
}

func init() {
	register(Operation{
		Name:     "ci.install",
		Summary:  "Install the Rhiza check as a GitHub workflow",
		Request:  reflect.TypeFor[CIInstallRequest](),
		Response: reflect.TypeFor[CIInstallResponse](),
		Run:      runCIInstall,
	})
}

// releaseVersion matches a module version go install can fetch: vX.Y.Z,
// optionally with a pre-release such as -rc.1 or a pseudo-version's
// -<time>-<hash>. Build metadata (+dirty from VCS stamping) never matches.
var releaseVersion = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$`)

// describeSuffix matches what git describe adds past a tag: -<n>-g<hash>.
var describeSuffix = regexp.MustCompile(`-[0-9]+-g[0-9a-f]+$`)

// installVersion is the version the workflow installs: v, when it names a
// published version, else "latest". A dirty build, a commits-past-a-tag
// git describe, a bare hash or "dev" name nothing go install can fetch.
func installVersion(v string) string {
	if !releaseVersion.MatchString(v) || strings.HasSuffix(v, "-dirty") || describeSuffix.MatchString(v) {
		return "latest"
	}
	return v
}

// withGitattribute returns have with the gitattribute line appended, after
// a newline if have lacks a final one, and true; or have itself and false
// when a line (less trailing spaces and \r) already equals it. Every other
// byte is kept.
func withGitattribute(have []byte) ([]byte, bool) {
	for l := range strings.SplitSeq(string(have), "\n") {
		if strings.TrimRight(l, " \r") == gitattribute {
			return have, false
		}
	}
	out := slices.Clone(have)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	return append(out, gitattribute+"\n"...), true
}

// runCIInstall writes the workflow and the .gitattributes line. A workflow
// on disk that differs from the template is a conflict unless req.Force;
// then nothing is written. Each file is written atomically, and only when
// its bytes change. ci.install writes no log entry, so --dry-run is a usage
// error, as for init.
func runCIInstall(ctx *Ctx, req any) (any, error) {
	r := req.(CIInstallRequest)
	if ctx.DryRun {
		return nil, Usagef("ci install does not take --dry-run: it writes no log entry").WithNext("rhiza ci install")
	}
	if _, err := openStore(ctx.Root); err != nil {
		return nil, err
	}
	resp := CIInstallResponse{Written: []string{}, Unchanged: []string{}}
	want := map[string][]byte{workflowPath: fmt.Appendf(nil, workflowTemplate, installVersion(ctx.Version))}

	wf, err := readOptional(ctx.Root, workflowPath)
	if err != nil {
		return nil, err
	}
	if wf != nil && !bytes.Equal(wf, want[workflowPath]) && !r.Force {
		return nil, Conflictf("%s exists and differs from the template", workflowPath).WithNext("rhiza ci install --force")
	}
	ga, err := readOptional(ctx.Root, gitattributesPath)
	if err != nil {
		return nil, err
	}
	want[gitattributesPath], _ = withGitattribute(ga)
	have := map[string][]byte{workflowPath: wf, gitattributesPath: ga}

	for _, p := range []string{gitattributesPath, workflowPath} {
		if have[p] != nil && bytes.Equal(have[p], want[p]) {
			resp.Unchanged = append(resp.Unchanged, p)
			continue
		}
		if err := writeAtomic(filepath.Join(ctx.Root, filepath.FromSlash(p)), want[p]); err != nil {
			return nil, Errorf("ci install: %v", err)
		}
		resp.Written = append(resp.Written, p)
	}
	return resp, nil
}

// readOptional reads the file at p under root; a missing file is nil. An
// empty file is a non-nil empty slice, so it reads as present.
func readOptional(root, p string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, Errorf("ci install: %v", err)
	}
	if b == nil {
		b = []byte{}
	}
	return b, nil
}

// Human reports the counts, then each path written.
func (r CIInstallResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Installed the Rhiza check (%d written, %d unchanged).\n", len(r.Written), len(r.Unchanged))
	if len(r.Written) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, p := range r.Written {
		fmt.Fprintf(w, "- wrote %s\n", p)
	}
}
