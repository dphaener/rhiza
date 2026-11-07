// Package dispatch reads a repo's .rhiza/dispatch.toml: the rules mapping a
// trigger and a request kind to the workflow a run starts.
package dispatch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// File is the rules file's path under a repo root.
const File = ".rhiza/dispatch.toml"

// The kinds a rule maps, and the tables that hold the rules.
const (
	KindFeature   = "feature"
	KindBug       = "bug"
	KindMilestone = "milestone"

	TableRequest = "request"
	TableReady   = "ready"
)

// The built-in defaults: today's behaviour.
const (
	DefaultRequest = "plan"
	DefaultReady   = "ship-change"
)

// Kinds lists the kinds a table may name.
var Kinds = []string{KindFeature, KindBug, KindMilestone}

// Rules maps each kind to a workflow name, per trigger. Every kind has an
// entry in both maps.
type Rules struct {
	Request map[string]string
	Ready   map[string]string
}

// Defaults returns the built-in rules.
func Defaults() Rules {
	r := Rules{Request: map[string]string{}, Ready: map[string]string{}}
	for _, k := range Kinds {
		r.Request[k] = DefaultRequest
		r.Ready[k] = DefaultReady
	}
	return r
}

// ForRequest is the workflow for a queued request of kind.
func (r Rules) ForRequest(kind string) string { return r.Request[kind] }

// ForReady is the workflow for an accepted, ready change set of kind.
func (r Rules) ForReady(kind string) string { return r.Ready[kind] }

// Load reads root's dispatch.toml. A missing file gives the defaults.
func Load(root string) (Rules, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(File)))
	if errors.Is(err, os.ErrNotExist) {
		return Defaults(), nil
	}
	if err != nil {
		return Rules{}, fmt.Errorf("read %s: %w", File, err)
	}
	r, err := Parse(data)
	if err != nil {
		return Rules{}, fmt.Errorf("%s: %w", File, err)
	}
	return r, nil
}

// Parse reads the file's content. Anything not [request] or [ready], any key
// that is not a kind, and any value that is not a bare workflow name is an
// error naming the table and key.
func Parse(data []byte) (Rules, error) {
	var raw map[string]any
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return Rules{}, err
	}
	r := Defaults()
	names := make([]string, 0, len(raw))
	for t := range raw {
		names = append(names, t)
	}
	sort.Strings(names)
	for _, t := range names {
		var dst map[string]string
		switch t {
		case TableRequest:
			dst = r.Request
		case TableReady:
			dst = r.Ready
		default:
			return Rules{}, fmt.Errorf("unknown table [%s]: want [%s] or [%s]", t, TableRequest, TableReady)
		}
		tbl, ok := raw[t].(map[string]any)
		if !ok {
			return Rules{}, fmt.Errorf("[%s] must be a table", t)
		}
		keys := make([]string, 0, len(tbl))
		for k := range tbl {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if !isKind(k) {
				return Rules{}, fmt.Errorf("[%s] unknown key %q: want one of %s", t, k, strings.Join(Kinds, ", "))
			}
			s, ok := tbl[k].(string)
			if !ok {
				return Rules{}, fmt.Errorf("[%s] %s: want a workflow name as a string, got %T", t, k, tbl[k])
			}
			if !bareName(s) {
				return Rules{}, fmt.Errorf("[%s] %s: %q is not a bare workflow name (no path, no .toml)", t, k, s)
			}
			dst[k] = s
		}
	}
	return r, nil
}

func isKind(k string) bool {
	for _, x := range Kinds {
		if x == k {
			return true
		}
	}
	return false
}

func bareName(s string) bool {
	return s != "" && !strings.ContainsAny(s, "/\\ \t\n") && !strings.HasSuffix(s, ".toml") && s != "." && s != ".."
}
