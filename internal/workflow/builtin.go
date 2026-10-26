package workflow

import (
	"embed"
	"path"
	"sort"
	"strings"
)

//go:embed builtin/*.toml
var builtinFiles embed.FS

// Builtin returns the source of the built-in workflow called name, and
// whether there is one.
func Builtin(name string) ([]byte, bool) {
	if name == "" || strings.ContainsAny(name, "/\\") {
		return nil, false
	}
	data, err := builtinFiles.ReadFile(path.Join("builtin", name+".toml"))
	if err != nil {
		return nil, false
	}
	return data, true
}

// BuiltinNames lists the built-in workflows' names, sorted.
func BuiltinNames() []string {
	entries, err := builtinFiles.ReadDir("builtin")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if n, ok := strings.CutSuffix(e.Name(), ".toml"); ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}
