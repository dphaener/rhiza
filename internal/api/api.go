// Package api defines Genesis's operations: every command, as a typed,
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

	"github.com/dphaener/genesis/internal/id"
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
	// Root is the project root: the folder holding .genesis.
	Root string
	// Home is the user's home directory, given by the caller (the API reads
	// no environment); empty if unknown. Operations that need it and find
	// it empty return a usage error.
	Home string
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
	// (the CLI's --change, --actor and --for, or GENESIS_CHANGE,
	// GENESIS_ACTOR and GENESIS_FOR), empty if none. Operations resolve
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
	// GitConfig returns the value of a git config key such as
	// "user.email", or an error if it is unset.
	GitConfig func(key string) (string, error)
	// Git runs git with args in Root and returns its standard output. It
	// fails when git is missing or exits non-zero. Nil means no git: the
	// operations that use it, such as check, skip what needs it.
	Git func(args ...string) (string, error)
	// Version is the running genesis's version, such as "v0.3.0" or
	// "dev".
	Version string
	// Stdin is read by operations that accept "-" for a file.
	Stdin io.Reader
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
					WithNext("genesis describe " + name)
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
