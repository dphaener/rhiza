// Command genesis is the thin CLI over the domain packages under internal/.
// It parses arguments, dispatches to a command and maps the result to an
// exit code; the behaviour itself lives in internal/.
package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// version is set at build time with -ldflags "-X main.version=...".
var version string

const usage = `Usage:
  genesis <command>

Commands:
  version    Print the genesis version
`

// versionString reports the ldflags version, else the module version from
// the build info, else "dev".
func versionString() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}
	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "genesis %s\n", versionString())
		return 0
	default:
		fmt.Fprintf(stderr, "genesis: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
