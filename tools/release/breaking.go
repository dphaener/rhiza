package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Surface is what a release promises: the API schemas, the daemon settings
// and the log format version, at one commit.
type Surface struct {
	Describe  []byte            // describe --json
	Settings  map[string]string // setting key -> type
	LogFormat int
}

// Break is one breaking difference between two surfaces.
type Break struct {
	Kind   string // operation, field, type, setting or log-format
	Detail string
}

func (b Break) String() string { return b.Kind + ": " + b.Detail }

type schema struct {
	Type       any               `json:"type"`
	Properties map[string]schema `json:"properties"`
	Items      *schema           `json:"items"`
}

type describeDoc struct {
	Operations []struct {
		Name     string `json:"name"`
		Request  schema `json:"request"`
		Response schema `json:"response"`
	} `json:"operations"`
}

// typeOf renders a schema's type, with an array's item type.
func typeOf(s schema) string {
	t := fmt.Sprint(s.Type)
	if b, err := json.Marshal(s.Type); err == nil {
		t = string(b)
	}
	if s.Items != nil {
		t += "[" + typeOf(*s.Items) + "]"
	}
	return t
}

// diffSchema reports removed and retyped fields of old under path, recursing
// into object properties and array items.
func diffSchema(path string, old, cur schema, out *[]Break) {
	if old.Items != nil && cur.Items != nil {
		diffSchema(path+"[]", *old.Items, *cur.Items, out)
	}
	names := make([]string, 0, len(old.Properties))
	for n := range old.Properties {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		op := old.Properties[n]
		np, ok := cur.Properties[n]
		if !ok {
			*out = append(*out, Break{"field", path + "." + n + " removed"})
			continue
		}
		if o, c := typeOf(op), typeOf(np); o != c {
			*out = append(*out, Break{"type", fmt.Sprintf("%s.%s changed from %s to %s", path, n, o, c)})
			continue
		}
		diffSchema(path+"."+n, op, np, out)
	}
}

// Breaking lists every breaking difference from old (the last release) to cur
// (HEAD): a removed operation or request/response field, a retyped field, a
// removed or retyped setting, a changed log format version. Additions are not
// breaking. The result is sorted by kind then detail, so it is stable.
func Breaking(old, cur Surface) ([]Break, error) {
	var o, c describeDoc
	if err := json.Unmarshal(old.Describe, &o); err != nil {
		return nil, fmt.Errorf("release describe: %w", err)
	}
	if err := json.Unmarshal(cur.Describe, &c); err != nil {
		return nil, fmt.Errorf("HEAD describe: %w", err)
	}
	var out []Break
	now := map[string]int{}
	for i, op := range c.Operations {
		now[op.Name] = i
	}
	for _, op := range o.Operations {
		i, ok := now[op.Name]
		if !ok {
			out = append(out, Break{"operation", op.Name + " removed"})
			continue
		}
		diffSchema(op.Name+" request", op.Request, c.Operations[i].Request, &out)
		diffSchema(op.Name+" response", op.Response, c.Operations[i].Response, &out)
	}
	for k, t := range old.Settings {
		nt, ok := cur.Settings[k]
		switch {
		case !ok:
			out = append(out, Break{"setting", k + " removed"})
		case nt != t:
			out = append(out, Break{"setting", fmt.Sprintf("%s changed type from %q to %q", k, t, nt)})
		}
	}
	if old.LogFormat != cur.LogFormat {
		out = append(out, Break{"log-format", fmt.Sprintf("version changed from %d to %d", old.LogFormat, cur.LogFormat)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Detail < out[j].Detail
	})
	return out, nil
}

var (
	settingHeadRE = regexp.MustCompile("(?m)^### `([^`]+)`$")
	typeLineRE    = regexp.MustCompile(`(?m)^- Type: (.+)$`)
	formatRE      = regexp.MustCompile(`(?m)^const FormatVersion = (\d+)`)
)

// parseSettings reads key -> type from docs/reference/config.md.
func parseSettings(md string) map[string]string {
	out := map[string]string{}
	parts := settingHeadRE.FindAllStringSubmatchIndex(md, -1)
	for i, p := range parts {
		end := len(md)
		if i+1 < len(parts) {
			end = parts[i+1][0]
		}
		t := typeLineRE.FindStringSubmatch(md[p[1]:end])
		if t != nil {
			out[md[p[2]:p[3]]] = strings.TrimSpace(t[1])
		}
	}
	return out
}

func parseLogFormat(src string) (int, error) {
	m := formatRE.FindStringSubmatch(src)
	if m == nil {
		return 0, fmt.Errorf("no FormatVersion in internal/entry/entry.go")
	}
	return strconv.Atoi(m[1])
}

func git(dir string, args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = dir
	var stderr bytes.Buffer
	c.Stderr = &stderr
	b, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(b), nil
}

// gitSurface reads the surface of ref in the repo at dir: it builds rhiza at
// ref in a throwaway worktree for describe --json, and reads the settings
// reference and the entry format from the tree.
func gitSurface(dir, ref string) (Surface, error) {
	var s Surface
	md, err := git(dir, "show", ref+":docs/reference/config.md")
	if err != nil {
		return s, err
	}
	s.Settings = parseSettings(md)
	src, err := git(dir, "show", ref+":internal/entry/entry.go")
	if err != nil {
		return s, err
	}
	if s.LogFormat, err = parseLogFormat(src); err != nil {
		return s, err
	}
	tmp, err := os.MkdirTemp("", "rhiza-breaking-")
	if err != nil {
		return s, err
	}
	defer os.RemoveAll(tmp)
	wt := filepath.Join(tmp, "wt")
	if _, err := git(dir, "worktree", "add", "--detach", wt, ref); err != nil {
		return s, err
	}
	defer git(dir, "worktree", "remove", "--force", wt)
	bin := filepath.Join(tmp, "rhiza")
	build := exec.Command("go", "build", "-trimpath", "-o", bin, "./cmd/rhiza")
	build.Dir = wt
	if out, err := build.CombinedOutput(); err != nil {
		return s, fmt.Errorf("build %s: %v\n%s", ref, err, out)
	}
	d := exec.Command(bin, "describe", "--json")
	d.Dir = tmp
	d.Env = append(os.Environ(), "RHIZA_ROOT=", "RHIZA_CHANGE=", "RHIZA_RUN=")
	if s.Describe, err = d.Output(); err != nil {
		return s, fmt.Errorf("describe at %s: %w", ref, err)
	}
	return s, nil
}

// breakingAgainstLastTag compares HEAD with the last v* tag in dir.
func breakingAgainstLastTag(dir string) (string, []Break, error) {
	tag, err := git(dir, "describe", "--tags", "--abbrev=0", "--match", "v*", "HEAD")
	if err != nil {
		return "", nil, fmt.Errorf("no v* tag to compare with: %w", err)
	}
	tag = strings.TrimSpace(tag)
	old, err := gitSurface(dir, tag)
	if err != nil {
		return tag, nil, err
	}
	cur, err := gitSurface(dir, "HEAD")
	if err != nil {
		return tag, nil, err
	}
	b, err := Breaking(old, cur)
	return tag, b, err
}
