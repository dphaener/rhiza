package api

// This file holds the init operation: it creates the repository, the
// .rhiza/.gitignore that keeps local/ (the current change set) out of git,
// and .rhiza/log/.gitkeep, since git keeps no empty folder and a clone of a
// fresh repository would otherwise have no log at all. It also asks which
// agent adapter the project uses and stores the answer in
// rhiza.project.agent.

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/store"
)

// gitignore is the content of .rhiza/.gitignore: local/ holds per-worktree
// state that must never be committed.
const gitignore = localDir + "/\n"

// InitRequest names the project's agent; init works in the caller's root.
type InitRequest struct {
	Agent string `json:"agent,omitempty" flag:"agent" desc:"The agent adapter the project uses, stored in rhiza.project.agent; default the answer to a prompt on a terminal, else the registered default adapter"`
}

// InitResponse names the new repository and the project's agent.
type InitResponse struct {
	Repo  id.ID  `json:"repo" desc:"The new repository's ID"`
	Root  string `json:"root" desc:"The absolute path of the project root holding .rhiza"`
	Agent string `json:"agent" desc:"The agent adapter the project uses"`
	Saved bool   `json:"saved" desc:"Whether the agent was written to rhiza.project.agent: false outside a git repository"`
}

func init() {
	register(Operation{
		Name:     "init",
		Summary:  "Create a Rhiza repository in the current folder",
		Request:  reflect.TypeFor[InitRequest](),
		Response: reflect.TypeFor[InitResponse](),
		Run:      runInit,
	})
}

// runInit creates the repository in ctx.Root. An existing one is a
// conflict. Init writes no log entry, so --dry-run has nothing to preview
// and is a usage error rather than a silent no-op. The log starts empty and
// git keeps no empty folder, so init writes an empty log/.gitkeep for the
// log to survive a clone; the log reader skips dot-files. The agent is
// chosen before anything is created, so an unknown --agent leaves no
// .rhiza; it is written to git config last, and only inside a git
// repository.
func runInit(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("init does not take --dry-run: it writes no log entry").WithNext("rhiza init")
	}
	root, err := ctx.abs(ctx.Root)
	if err != nil {
		return nil, Errorf("project root: %v", err)
	}
	agent, err := initAgent(ctx, root, req.(InitRequest).Agent)
	if err != nil {
		return nil, err
	}
	st, err := store.Init(root, ctx.IDs)
	if err != nil {
		if errors.Is(err, store.ErrOldFolder) {
			return nil, oldFolderError()
		}
		if errors.Is(err, store.ErrExists) {
			return nil, Conflictf("a Rhiza repository already exists in %s", root).WithNext("rhiza change ls")
		}
		return nil, Errorf("init: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, store.Dir, ".gitignore"), []byte(gitignore), 0o644); err != nil {
		return nil, Errorf("write .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, store.Dir, store.LogDir, ".gitkeep"), nil, 0o644); err != nil {
		return nil, Errorf("write log/.gitkeep: %v", err)
	}
	resp := InitResponse{Repo: st.Repo().ID, Root: root, Agent: agent}
	if ctx.Git != nil {
		if _, err := ctx.Git("rev-parse", "--git-dir"); err == nil {
			if _, err := ctx.Git("config", ProjectAgentKey, agent); err != nil {
				return nil, Errorf("write %s: %v", ProjectAgentKey, err)
			}
			resp.Saved = true
		}
	}
	return resp, nil
}

// initAgent is the project's agent: the request's, which must be
// registered; else, on a terminal and not --json, the answer to a prompt;
// else the registered default. Where a repository already exists init
// fails as a conflict, so it asks nothing there.
func initAgent(ctx *Ctx, root, flag string) (string, error) {
	if flag = strings.TrimSpace(flag); flag != "" {
		if err := validateAgent(flag); err != nil {
			return "", Usagef("--agent %q %v", flag, err).WithNext("rhiza init --agent <agent>")
		}
		return flag, nil
	}
	d := adapter.Default()
	if d == nil {
		return "", Errorf("no default agent adapter is registered")
	}
	if ctx.Terminal == nil || ctx.Stdin == nil || ctx.JSON {
		return d.Name(), nil
	}
	if _, found := FindRoot(root); found {
		return d.Name(), nil
	}
	return askAgentName(ctx.Terminal, ctx.Stdin, d.Name())
}

// askAgentName prints the registered adapters to term, def marked as the
// default, and reads one line from in: an empty line takes def, an unknown
// name asks again. Input that ends first takes def when nothing was typed.
func askAgentName(term io.Writer, in io.Reader, def string) (string, error) {
	fmt.Fprintln(term, "Which agent does this project use?")
	for _, n := range adapter.Names() {
		if n == def {
			fmt.Fprintf(term, "  %s (default)\n", n)
		} else {
			fmt.Fprintf(term, "  %s\n", n)
		}
	}
	r := bufio.NewReader(in)
	for {
		fmt.Fprintf(term, "Agent [%s]: ", def)
		line, err := r.ReadString('\n')
		name := strings.TrimSpace(line)
		if name == "" {
			if err != nil {
				fmt.Fprintln(term)
			}
			return def, nil
		}
		verr := validateAgent(name)
		if verr == nil {
			return name, nil
		}
		fmt.Fprintf(term, "%s %v\n", name, verr)
		if err != nil {
			return "", Usagef("agent %q %v", name, verr).WithNext("rhiza init --agent <agent>")
		}
	}
}

// Human reports the new repository, the project's agent and the first
// command to run.
func (r InitResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Initialised a Rhiza repository in %s (repo %s).\n", r.Root, r.Repo)
	if r.Saved {
		fmt.Fprintf(w, "Agent: %s, written to %s in git config.\n\n", r.Agent, ProjectAgentKey)
	} else {
		fmt.Fprintf(w, "Agent: %s, not written: this is not a git repository. Once it is one, run rhiza daemon set %s %s.\n\n", r.Agent, ProjectAgentKey, r.Agent)
	}
	fmt.Fprintf(w, "next: %s\n", newEntryNext)
}
