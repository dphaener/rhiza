// Package api defines Rhiza's operations: every command, as a typed,
// transport-independent call with a published schema. The CLI and any future
// server are only callers; they turn their input into a request, run the
// operation and render its response.
//
// The vocabulary:
//
//   - An operation is one named call, such as "init", "change.new" or
//     "add". Names mirror command paths, with dots between the words.
//   - A request is the operation's input, a Go struct decoded strictly from
//     JSON (see DecodeRequest). A response is its output, a Go struct
//     encoded as canonical JSON.
//   - A schema is the JSON Schema of a request or response type, reflected
//     from its Go struct tags (see Schema).
//   - A Ctx carries everything a call needs that is not in its request: the
//     project root, the change set, the actor, the clock and so on. Callers
//     fill it explicitly; nothing is read from the environment here.
//   - An Error is a failure with a code, a message and an optional next
//     command; its code decides the process exit code (see ExitCode).
//
// Operations register themselves from an init function in their own file,
// so adding one never edits a shared list.
package api

import (
	"cmp"
	"io"
	"reflect"
	"slices"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/store"
)

// Operation is one named call. Request and Response are struct types; Run
// receives a value of type Request (not a pointer, as DecodeRequest returns
// it) and returns a value of type Response.
type Operation struct {
	// Name is the dotted command path, such as "change.new".
	Name string
	// Summary is one line describing the operation.
	Summary string
	// Write reports whether the operation appends to the log.
	Write bool
	// ReadOnly marks an operation that changes nothing: not the log, the
	// settings, any file, Rex or GitHub. The web UI's /api/<op> runs only
	// operations with it. Write does not stand for it, since it only means
	// appending to the log.
	ReadOnly bool
	// Record reports whether the operation takes Ctx.Record. Run refuses
	// Ctx.Record with a usage error for an operation without it.
	Record bool
	// Request and Response are the types of the call's input and output.
	Request, Response reflect.Type
	// Input, if not nil, describes a file or stdin document the operation
	// reads besides its request, such as apply's file of writes.
	Input *Input
	// Run performs the call.
	Run func(ctx *Ctx, req any) (any, error)
}

// Streamer is implemented by a response that went to Ctx.Out as it
// happened: when Streamed reports true, callers print nothing for it.
type Streamer interface {
	Streamed() bool
}

// Exiter is implemented by a response that sets the process exit code, such
// as a check that failed. Callers print the response in full, as markdown or
// JSON as usual, and then exit with ExitCode.
type Exiter interface {
	ExitCode() int
}

// Input describes a document an operation reads from a file or stdin.
type Input struct {
	// Summary is one line describing the document's format for people.
	Summary string
	// Schema returns the document's JSON Schema, ready for canon.Marshal.
	// It is called lazily, so it may look up other operations.
	Schema func() map[string]any
}

// Ctx carries everything a call needs besides its request. Callers set every
// field explicitly; the CLI fills its defaults from flags, the environment
// and the working tree.
//
// The change set, actor and for-actor may each be given as an ID or as a
// raw ref, not both (that is a usage error). An ID wins over a ref, and a
// ref over the default: the current change set, and the actor whose email
// is git's user.email (recorded on first use). The CLI fills only the refs;
// a programmatic caller, such as a server with no current anything, may set
// the IDs instead.
type Ctx struct {
	// Wd is the caller's working directory, absolute; empty if unknown.
	// rhiza root init and root show start from it.
	Wd string
	// Root is the project root: the folder holding .rhiza.
	Root string
	// Home is the user's home directory, given by the caller (the API reads
	// no environment); empty if unknown. Operations that need it and find
	// it empty return a usage error.
	Home string
	// Getenv looks up a variable of the caller's environment, for adapters
	// that locate an agent's own files, such as a session's transcript
	// (Home answers HOME when set). Nil means no environment: an operation
	// that needs it finds nothing there, and session.run captures no
	// answers when Home is empty too.
	Getenv func(key string) string
	// IndexPath is the index database file for Root, given by the caller,
	// since nothing here reads the environment or the user cache dir (the
	// CLI uses index.DefaultPath). Empty is an error from any read.
	IndexPath string
	// ChangeSet is the change set to act on, zero to choose by ChangeRef
	// or the current one. It must be in the log, else not_found.
	ChangeSet id.ID
	// Actor is who performs the call, For whom it acts for; zero to choose
	// by ActorRef and ForRef or the default. Each must be a live actor node.
	Actor, For id.ID
	// ChangeRef, ActorRef and ForRef are the raw refs the caller was given
	// (the CLI's --change, --actor and --for, or RHIZA_CHANGE,
	// RHIZA_ACTOR and RHIZA_FOR), empty if none. Operations resolve
	// them when the matching ID is zero.
	ChangeRef, ActorRef, ForRef string
	// Record asks a write to open a record change set, write into it and
	// land it in one step, leaving the current change set alone. Only
	// operations with Operation.Record take it.
	Record bool
	// DryRun asks a write to build everything and write nothing.
	DryRun bool
	// Now is the clock.
	Now func() time.Time
	// IDs makes new IDs.
	IDs *id.Generator
	// readLog, when set, replaces the store's log read; tests count folds
	// with it.
	readLog func() store.Log
	// builtModel, when set, is told each change set whose model an index
	// read builds (memo misses only); tests count folds per change set
	// with it.
	builtModel func(cs id.ID)
	// daemonFold is the daemon's fold cache; the first poll creates it.
	daemonFold *daemonFoldCache
	// GitConfig returns the value of a git config key such as
	// "user.email", or an error if it is unset.
	GitConfig func(key string) (string, error)
	// Git runs git with args in Root and returns its standard output. It
	// fails when git is missing or exits non-zero. Nil means no git: the
	// operations that use it, such as check, skip what needs it.
	Git func(args ...string) (string, error)
	// Exec runs the program name with args in dir and returns its standard
	// output. A nil env means the caller's environment. A failure carries
	// the program's trimmed standard error. Nil means no programs can be
	// run: the operations that need one, such as session.start, fail.
	Exec func(name string, args []string, dir string, env []string) ([]byte, error)
	// LookPath finds a program on the caller's PATH, as exec.LookPath does.
	// Nil means every program is missing: daemon.install finds no systemd.
	LookPath func(file string) (string, error)
	// Self is the absolute path of the running rhiza, for commands that
	// must call it back, such as the command a session's terminal tab runs.
	// Empty if unknown.
	Self string
	// Start starts a program and returns at once; the caller waits for it.
	// Unlike Exec it can hand the program the caller's terminal, or stream
	// its output. Nil means no programs can be started: session.run fails.
	Start func(p Program) (Process, error)
	// Detach starts p as a child that outlives the caller: in its own
	// session, with no input, its output appended to the file logPath
	// (created with its folders). It returns once the child has started and
	// does not wait for it. Nil means no program can be detached: the
	// daemon cannot start a run.
	Detach func(p Program, logPath string) error
	// Reap waits, without blocking, on the detached children of the repo
	// whose git common dir is common that have exited and that this process
	// still parents, so none stays a zombie. Nil means nothing to reap.
	Reap func(common string)
	// RexSession is the Rex session the caller runs in (the CLI's
	// REX_SESSION), empty if none: session.start opens its tab there
	// unless the request names another.
	RexSession string
	// repoSession is set once the daemon has made sure its Rex session
	// rhiza-<repo> exists and made it RexSession.
	repoSession bool
	// ProjectName is the name the project is registered under in the
	// root's rhiza.toml, set by the root daemon; empty elsewhere.
	ProjectName string
	// RexBlock is the Rex block the caller runs in (the CLI's REX_BLOCK),
	// empty if none: session.run copies that block's status reports.
	RexBlock string
	// Stderr receives diagnostics an operation must not print as its
	// response, such as a failed hook's; nil discards them.
	Stderr io.Writer
	// Rex is the terminal Rex shows this process in, nil outside Rex.
	// run.start and run.resume write their OSC 7501 status reports to it,
	// such as blocked while an agent step waits on the owner; nothing is
	// ever read back from it.
	Rex io.Writer
	// Version is the running rhiza's version, such as "v0.3.0" or
	// "dev".
	Version string
	// Commit is the git commit the running rhiza was built from, or "".
	Commit string
	// Stdin is read by operations that accept "-" for a file.
	Stdin io.Reader
	// Out is the sink an operation that streams (feed --follow) writes its
	// lines to as they happen; nil means no streaming. JSON asks such an
	// operation for NDJSON instead of text lines.
	Out  io.Writer
	JSON bool
	// Interrupts starts listening for the caller being asked to stop (the
	// CLI's interrupt and terminate signals): the channel it returns gets
	// the first one, and stop ends the listening. An operation that runs
	// until stopped calls it, so the others keep the default handling of
	// the signals. Nil means the caller cannot be asked to stop.
	Interrupts func() (interrupted <-chan struct{}, stop func())
}

// Program is a program for Ctx.Start.
type Program struct {
	// Name is the program, looked up on PATH unless it has a separator.
	Name string
	// Args are its arguments, without the name.
	Args []string
	// Dir is the folder it runs in; empty means the caller's.
	Dir string
	// Env holds NAME=value pairs set over the caller's environment.
	Env []string
	// Unset names variables removed from the caller's environment.
	Unset []string
	// Terminal hands the program the caller's standard input, output and
	// error, and makes the caller ignore the terminal's interrupts while
	// it runs, so they reach the program alone. Without it the program
	// reads nothing, its standard output is Process.Output and its
	// standard error is discarded.
	Terminal bool
}

// Process is a started Program.
type Process interface {
	// Output is the program's standard output; nil for a Terminal program.
	// Read it to the end before calling Wait.
	Output() io.Reader
	// Wait waits for the program to end and returns its exit code: 128
	// plus the signal number if a signal ended it. The error is set only
	// when the program could not be waited for.
	Wait() (int, error)
	// Kill ends the program. Output then reaches its end, at the latest
	// shortly after, even if a child of the program still holds it.
	Kill()
}

// registry holds operations by name. Tests use their own; the package uses
// defaultRegistry.
type registry struct {
	ops map[string]Operation
}

var defaultRegistry = &registry{}

// add registers o, wrapping Run to refuse Ctx.Record unless o.Record is
// set. It panics if o has no name, a name already registered,
// or no Request, Response or Run, since these are programming errors.
func (r *registry) add(o Operation) {
	switch {
	case o.Name == "":
		panic("api: register: empty operation name")
	case o.Request == nil:
		panic("api: register " + o.Name + ": nil Request type")
	case o.Response == nil:
		panic("api: register " + o.Name + ": nil Response type")
	case o.Run == nil:
		panic("api: register " + o.Name + ": nil Run")
	}
	if _, ok := r.ops[o.Name]; ok {
		panic("api: register " + o.Name + ": duplicate operation name")
	}
	if r.ops == nil {
		r.ops = map[string]Operation{}
	}
	if !o.Record {
		run, name := o.Run, o.Name
		o.Run = func(ctx *Ctx, req any) (any, error) {
			if ctx != nil && ctx.Record {
				return nil, Usagef("%s does not take --record: it is for node and link writes", CommandPath(name)).
					WithNext("rhiza describe " + name)
			}
			return run(ctx, req)
		}
	}
	r.ops[o.Name] = o
}

func (r *registry) all() []Operation {
	out := make([]Operation, 0, len(r.ops))
	for _, o := range r.ops {
		out = append(out, o)
	}
	slices.SortFunc(out, func(a, b Operation) int { return cmp.Compare(a.Name, b.Name) })
	return out
}

func (r *registry) lookup(name string) (Operation, bool) {
	o, ok := r.ops[name]
	return o, ok
}

// register adds o to the package's operations. Call it from an init
// function; it panics on an empty or duplicate name.
func register(o Operation) {
	defaultRegistry.add(o)
}

// Operations returns every registered operation, sorted by name.
func Operations() []Operation {
	return defaultRegistry.all()
}

// Lookup returns the operation with the given name.
func Lookup(name string) (Operation, bool) {
	return defaultRegistry.lookup(name)
}
