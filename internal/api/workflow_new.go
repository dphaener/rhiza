package api

// This file holds workflow.new: a workflow written into the repo's
// .rhiza/workflows folder, forked from a resolved workflow or started from a
// skeleton. It never overwrites a file.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "workflow.new",
		Summary:  "Write .rhiza/workflows/<name>.toml: a copy of a resolved workflow with its name set (--from), else a minimal commented skeleton; refuses an existing file and validates the result as workflow check does",
		Request:  reflect.TypeFor[WorkflowNewRequest](),
		Response: reflect.TypeFor[WorkflowNewResponse](),
		Run:      runWorkflowNew,
	})
}

// WorkflowNewRequest names the new workflow and, optionally, its source.
type WorkflowNewRequest struct {
	Name string `json:"name" arg:"0" desc:"The new workflow's bare name (no / and no .toml)"`
	From string `json:"from,omitempty" flag:"from" desc:"A workflow name to fork (personal, repo, then built-in) or a .toml file path; without it a skeleton is written"`
}

// WorkflowNewResponse is the file written, or in a dry run the file that would be.
type WorkflowNewResponse struct {
	Name   string `json:"name" desc:"The new workflow's name"`
	Path   string `json:"path" desc:"The file written"`
	From   string `json:"from,omitempty" desc:"Where the copy came from; empty for a skeleton"`
	DryRun bool   `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// Human writes the path.
func (r WorkflowNewResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "would write %s\n", r.Path)
		return
	}
	fmt.Fprintln(w, r.Path)
}

const workflowSkeleton = `name = "%[1]s"

# A workflow is steps run in order. Check it with
#   rhiza workflow check %[1]s
# and see what a run would do with
#   rhiza run start --workflow %[1]s --dry-run
start = "build"

# An agent step runs an agent session with this prompt. {change}, {title},
# {step} and {run} are filled in.
[step.build]
kind = "agent"
permission_mode = "manual"
agent = "claude-code"
prompt = "Build change set {change} ({title}) in step {step} of run {run}"
on.pass = "test"

# A command step runs a shell command; on.pass and on.fail name the next
# step, or done or failed.
[step.test]
kind = "command"
run = "make test"
max_visits = 3
on.pass = "done"
on.fail = "build"
`

var (
	topLevelName = regexp.MustCompile(`(?m)^name[ \t]*=.*$`)
	firstTable   = regexp.MustCompile(`(?m)^\[`)
)

// withName returns src with its top-level name set to name and everything
// else byte for byte as it was.
func withName(src []byte, name string) []byte {
	s := string(src)
	line := "name = " + strconv.Quote(name)
	head, rest := s, ""
	if loc := firstTable.FindStringIndex(s); loc != nil {
		head, rest = s[:loc[0]], s[loc[0]:]
	}
	if loc := topLevelName.FindStringIndex(head); loc != nil {
		return []byte(head[:loc[0]] + line + head[loc[1]:] + rest)
	}
	return []byte(line + "\n\n" + head + rest)
}

func runWorkflowNew(ctx *Ctx, req any) (any, error) {
	r := req.(WorkflowNewRequest)
	name := strings.TrimSpace(r.Name)
	const next = "rhiza workflow new <name> [--from <workflow>]"
	switch {
	case name == "":
		return nil, Usagef("workflow new needs a name").WithNext(next)
	case name == "." || name == ".." || isWorkflowPath(name):
		return nil, Usagef("%q is not a workflow name: give a bare name, with no / and no .toml", name).WithNext(next)
	}
	var src []byte
	from := strings.TrimSpace(r.From)
	if from == "" {
		src = fmt.Appendf(nil, workflowSkeleton, name)
	} else {
		var data []byte
		if isWorkflowPath(from) {
			p, err := workflowFile(ctx.Root, from)
			if err != nil {
				return nil, err
			}
			if data, err = os.ReadFile(p); err != nil {
				return nil, Usagef("no workflow %q: %v", from, err).WithNext(next)
			}
			from = p
		} else {
			p, d, ok, err := lookupWorkflow(ctx.Root, from)
			if err != nil {
				return nil, Usagef("workflow %q: %v", from, err).WithNext(next)
			}
			if !ok {
				return nil, Usagef("no workflow %q to copy: there is no such file in %s or your personal workflows, and no built-in of that name (built-in: %s)",
					from, workflowsDir, strings.Join(workflow.BuiltinNames(), ", ")).WithNext(next)
			}
			data, from = d, p
		}
		src = withName(data, name)
	}
	path := filepath.Join(ctx.Root, filepath.FromSlash(workflowsDir), name+".toml")
	if _, err := os.Lstat(path); err == nil {
		return nil, Conflictf("%s already exists; it is left as it is", path).WithNext("rhiza workflow describe " + name)
	}
	if _, problems := workflow.Validate(name, src, workflowResolver(ctx.Root)); len(problems) > 0 {
		var b strings.Builder
		for _, p := range problems {
			if p.Step != "" {
				fmt.Fprintf(&b, "\n  step %s: %s", p.Step, p.Message)
			} else {
				fmt.Fprintf(&b, "\n  %s", p.Message)
			}
		}
		return nil, Usagef("the new workflow %s would not pass workflow check:%s", name, b.String())
	}
	out := WorkflowNewResponse{Name: name, Path: path, From: from, DryRun: ctx.DryRun}
	if ctx.DryRun {
		return out, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, Errorf("workflow new: %v", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return nil, Conflictf("%s already exists; it is left as it is", path)
	}
	if err != nil {
		return nil, Errorf("workflow new: %v", err)
	}
	if _, err := f.Write(src); err != nil {
		f.Close()
		return nil, Errorf("workflow new: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, Errorf("workflow new: %v", err)
	}
	return out, nil
}
