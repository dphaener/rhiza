// Command rhiza is the thin CLI over the domain packages under internal/.
// It is a generic frame over the operations in internal/api: it maps argv
// to a request through the request type's struct tags, runs the operation,
// renders the response as markdown or JSON and maps errors to exit codes.
// It has no per-command code; only version, api call, tui and help are its own.
package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"strings"

	"github.com/dphaener/rhiza/internal/api"
)

// version is set at build time with -ldflags "-X main.version=...".
var version string

// versionString reports the ldflags version (the release tag, e.g. v0.1.0),
// else the module version from the build info (go install), else "dev".
func versionString() string {
	mod := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		mod = info.Main.Version
	}
	return pickVersion(version, mod)
}

// pickVersion prefers the ldflags version, then a real module version, and
// falls back to "dev" (a "(devel)" module version is a local build).
func pickVersion(ldflags, module string) string {
	if ldflags != "" {
		return ldflags
	}
	if module != "" && module != "(devel)" {
		return module
	}
	return "dev"
}

// commitString is the git commit the binary was built from (the build
// info's vcs.revision), or "". RHIZA_TEST_COMMIT stands in for it
// alongside RHIZA_TEST_NOW, since a test binary carries none.
func commitString(getenv func(string) string) string {
	if v := getenv("RHIZA_TEST_COMMIT"); v != "" && getenv("RHIZA_TEST_NOW") != "" {
		return v
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				return s.Value
			}
		}
	}
	return ""
}

const globalFlags = `Global flags:
  --json          Print JSON instead of markdown
  --dry-run       Build a write and print it; write nothing
  --change <ref>  The change set to act on (env RHIZA_CHANGE)
  --actor <ref>   Who performs the call (env RHIZA_ACTOR)
  --for <ref>     Whom the call acts for (env RHIZA_FOR)
  --record        Write in a record change set landed in one step
  -h, --help      Show this help, or describe a command
`

// usage is the help text, generated from the registry.
func usage() string {
	type row struct{ name, summary string }
	var rows []row
	for _, o := range api.Operations() {
		rows = append(rows, row{api.CommandPath(o.Name), o.Summary})
	}
	for _, f := range api.FrameCommands {
		rows = append(rows, row{f.Name, f.Summary})
	}
	width := 0
	for _, r := range rows {
		width = max(width, len(r.name))
	}
	var b strings.Builder
	b.WriteString("Usage:\n  rhiza <command> [arguments] [flags]\n  rhiza api call <operation> [<json>|-]\n\nCommands:\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "  %-*s  %s\n", width, r.name, r.summary)
	}
	b.WriteString("\n" + globalFlags + "\nRun \"rhiza describe <operation>\" for an operation's arguments.\n")
	return b.String()
}

// run is the testable entry point: it returns the exit code.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) int {
	g, rest, err := parseGlobals(args)
	if err != nil {
		return writeError(stdout, stderr, err, g.json)
	}
	words := commandWords(rest)
	if len(words) == 0 {
		if len(rest) > 0 && !g.help {
			return writeError(stdout, stderr, api.Usagef("unknown flag %s", rest[0]).WithNext("rhiza --help"), g.json)
		}
		fmt.Fprint(stdout, usage())
		return 0
	}
	switch words[0] {
	case "version":
		if g.help {
			fmt.Fprint(stdout, usage())
			return 0
		}
		fmt.Fprintf(stdout, "rhiza %s\n", versionString())
		return 0
	case "tui":
		if g.help {
			fmt.Fprint(stdout, usage())
			return 0
		}
		return tuiCommand(g, rest, stdout, stderr, getenv)
	case "api":
		if g.help {
			fmt.Fprint(stdout, usage())
			return 0
		}
		return apiCall(g, rest[1:], stdin, stdout, stderr, getenv)
	}
	o, n, ok := matchCommand(words)
	if !ok {
		if path, subs, ok := subcommands(words); ok {
			return writeError(stdout, stderr, api.Usagef("%q needs a subcommand: %s", path, strings.Join(subs, ", ")).WithNext("rhiza --help"), g.json)
		}
		if g.json {
			return writeError(stdout, stderr, api.Usagef("unknown command %q", words[0]).WithNext("rhiza --help"), true)
		}
		fmt.Fprintf(stderr, "rhiza: unknown command %q\n\n%s", words[0], usage())
		return 2
	}
	if g.help {
		// Help describes the command; --record is for the command, not
		// for describe.
		g.record = false
		return call(g, describeOp(), api.DescribeRequest{Op: o.Name}, stdin, stdout, stderr, getenv)
	}
	req, err := requestFromArgs(o, rest[n:])
	if err != nil {
		return writeError(stdout, stderr, err, g.json)
	}
	return call(g, o, req, stdin, stdout, stderr, getenv)
}

// apiCall runs `rhiza api call <op> [<json>|-]`. Its output is always
// JSON.
func apiCall(g globals, args []string, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) int {
	g.json = true
	if len(args) == 0 || args[0] != "call" {
		return writeError(stdout, stderr, api.Usagef("usage: rhiza api call <operation> [<json>|-]").WithNext("rhiza describe"), true)
	}
	args = args[1:]
	for _, a := range args {
		if a != "-" && strings.HasPrefix(a, "-") {
			return writeError(stdout, stderr, api.Usagef("api call: unknown flag %s", a), true)
		}
	}
	if len(args) < 1 || len(args) > 2 {
		return writeError(stdout, stderr, api.Usagef("usage: rhiza api call <operation> [<json>|-]").WithNext("rhiza describe"), true)
	}
	o, ok := api.Lookup(args[0])
	if !ok {
		return writeError(stdout, stderr, api.NotFoundf("unknown operation %q", args[0]).WithNext("rhiza describe"), true)
	}
	var data []byte
	if len(args) == 2 {
		data = []byte(args[1])
		if args[1] == "-" {
			b, err := io.ReadAll(stdin)
			if err != nil {
				return writeError(stdout, stderr, fmt.Errorf("read stdin: %w", err), true)
			}
			data = b
		}
	}
	req, err := api.DecodeRequest(o, data)
	if err != nil {
		return writeError(stdout, stderr, err, true)
	}
	return call(g, o, req, stdin, stdout, stderr, getenv)
}

// call runs o on req and prints the outcome. A response that is an
// api.Exiter sets the exit code once it is printed.
func call(g globals, o api.Operation, req any, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) int {
	ctx, err := newCtx(g, stdin, getenv)
	if err != nil {
		return writeError(stdout, stderr, err, g.json)
	}
	ctx.Out, ctx.JSON = stdout, g.json
	resp, err := o.Run(ctx, req)
	if s, ok := resp.(api.Streamer); ok && err == nil && s.Streamed() {
		return 0
	}
	if err == nil {
		err = writeResponse(stdout, resp, g.json)
	}
	if err != nil {
		return writeError(stdout, stderr, err, g.json)
	}
	if e, ok := resp.(api.Exiter); ok {
		return e.ExitCode()
	}
	return 0
}

func describeOp() api.Operation {
	o, ok := api.Lookup("describe")
	if !ok {
		panic("rhiza: no describe operation")
	}
	return o
}

func main() {
	os.Exit(profiled(os.Getenv("RHIZA_CPUPROFILE"), func() int {
		return run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Getenv)
	}))
}

// profiled runs f, writing a CPU profile to path when path is set. It is
// a hidden diagnostic: RHIZA_CPUPROFILE=<file> rhiza <cmd>, then
// go tool pprof -top <file>.
func profiled(path string, f func() int) int {
	if path == "" {
		return f()
	}
	out, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rhiza: cpu profile: %v\n", err)
		return f()
	}
	defer out.Close()
	if err := pprof.StartCPUProfile(out); err != nil {
		fmt.Fprintf(os.Stderr, "rhiza: cpu profile: %v\n", err)
		return f()
	}
	defer pprof.StopCPUProfile()
	return f()
}
