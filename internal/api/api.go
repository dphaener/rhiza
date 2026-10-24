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
	// Request and Response are the types of the call's input and output.
	Request, Response reflect.Type
	// Run performs the call.
	Run func(ctx *Ctx, req any) (any, error)
}

// Ctx carries everything a call needs besides its request. Callers set every
// field explicitly; the CLI fills its defaults from flags, the environment
// and the working tree.
type Ctx struct {
	// Root is the project root: the folder holding .genesis.
	Root string
	// ChangeSet is the change set writes go into; zero if there is none.
	ChangeSet id.ID
	// Actor is who performs the call, For whom it acts for. Either may be
	// zero when the caller leaves the choice to the operation.
	Actor, For id.ID
	// DryRun asks a write to build everything and write nothing.
	DryRun bool
	// Now is the clock.
	Now func() time.Time
	// IDs makes new IDs.
	IDs *id.Generator
	// GitConfig returns the value of a git config key such as
	// "user.email", or an error if it is unset.
	GitConfig func(key string) (string, error)
	// Stdin is read by operations that accept "-" for a file.
	Stdin io.Reader
}

// registry holds operations by name. Tests use their own; the package uses
// defaultRegistry.
type registry struct {
	ops map[string]Operation
}

var defaultRegistry = &registry{}

// add registers o. It panics if o has no name, a name already registered,
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
