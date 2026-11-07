// Package rootcfg is the Rhiza root: a directory holding a rhiza.toml, which
// registers the projects under it and holds machine-wide settings. Every
// rhiza command finds its root by walking up from the working directory to
// the nearest rhiza.toml.
package rootcfg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// File is the name of a root's config file.
const File = "rhiza.toml"

// ErrExists is returned by Init when the directory is already a root.
var ErrExists = errors.New("already a Rhiza root")

// Project is one registered project.
type Project struct {
	// Path is the project's folder, relative to the root (or absolute).
	Path string `toml:"path"`
	// Repo is the project's repo ID, when it has .rhiza/.
	Repo string `toml:"repo,omitempty"`
}

// Config is the content of a root's rhiza.toml.
type Config struct {
	// SelfUpdateProject names the registered project whose main checkout
	// the daemon keeps itself current from (self_update_project). Empty
	// means the daemon never self-updates.
	SelfUpdateProject string             `toml:"self_update_project,omitempty"`
	Projects          map[string]Project `toml:"projects"`
	// Daemon holds the machine-wide daemon settings, by their full key
	// (rhiza.daemon.maxRuns), each a string as rhiza daemon set takes it.
	Daemon map[string]string `toml:"daemon,omitempty"`
}

// Find walks up from dir to the nearest folder holding a rhiza.toml and
// returns that folder. It reports false if no ancestor has one, and for an
// empty dir: that is no directory at all, never the process's cwd.
func Find(dir string) (string, bool) {
	if dir == "" {
		return "", false
	}
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	for {
		if fi, err := os.Stat(filepath.Join(d, File)); err == nil && !fi.IsDir() {
			return d, true
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}

// Init makes dir (created if missing) a root by writing an empty
// rhiza.toml, and returns the config file's path. It returns ErrExists,
// writing nothing, if dir already holds one.
func Init(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(abs, File)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return path, ErrExists
		}
		return "", err
	}
	_, werr := f.WriteString("# The Rhiza root: its projects and machine-wide settings.\n")
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		os.Remove(path)
		return "", werr
	}
	return path, nil
}

// Load reads the root's rhiza.toml.
func Load(root string) (Config, error) {
	var c Config
	path := filepath.Join(root, File)
	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if _, err := toml.Decode(string(data), &c); err != nil {
		return c, fmt.Errorf("%s: %w", path, err)
	}
	return c, nil
}

// Names returns the registered project names, sorted.
func (c Config) Names() []string {
	var out []string
	for n := range c.Projects {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Dir is the absolute folder of the project name registered under root.
func (c Config) Dir(root, name string) string {
	p := c.Projects[name].Path
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	return filepath.Join(root, p)
}

var slugRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// ProjectOf names the project that dir belongs to: the registered project
// whose folder is dir or holds it, or whose <project>.<slug> sibling
// worktree is. It reports false at the root itself or anywhere else.
func (c Config) ProjectOf(root, dir string) (string, bool) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	for _, name := range c.Names() {
		pd := c.Dir(root, name)
		if within(pd, d) {
			return name, true
		}
		// A sibling worktree: <project folder>.<slug>, directly beside it.
		base := filepath.Base(pd) + "."
		rel, err := filepath.Rel(filepath.Dir(pd), d)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			continue
		}
		first := strings.SplitN(filepath.ToSlash(rel), "/", 2)[0]
		if slug, ok := strings.CutPrefix(first, base); ok && slugRE.MatchString(slug) {
			return name, true
		}
	}
	return "", false
}

func within(parent, d string) bool {
	rel, err := filepath.Rel(parent, d)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// ErrTaken is returned by Add when the name is already registered.
var ErrTaken = errors.New("project name already registered")

// Add registers project name at path (relative to root) in root's
// rhiza.toml, keeping the other projects. The file is rewritten whole, so
// comments in it are not kept except the header.
func Add(root, name string, p Project) error {
	c, err := Load(root)
	if err != nil {
		return err
	}
	if _, ok := c.Projects[name]; ok {
		return ErrTaken
	}
	if c.Projects == nil {
		c.Projects = map[string]Project{}
	}
	c.Projects[name] = p
	return write(root, c)
}

// ErrNotRegistered is returned by Remove when the name is not registered.
var ErrNotRegistered = errors.New("project not registered")

// Remove unregisters project name from root's rhiza.toml, keeping the other
// projects. It never touches the project's files.
func Remove(root, name string) error {
	c, err := Load(root)
	if err != nil {
		return err
	}
	if _, ok := c.Projects[name]; !ok {
		return ErrNotRegistered
	}
	delete(c.Projects, name)
	return write(root, c)
}

// SetDaemon writes the machine setting key = value in root's rhiza.toml,
// keeping the rest.
func SetDaemon(root, key, value string) error {
	c, err := Load(root)
	if err != nil {
		return err
	}
	if c.Daemon == nil {
		c.Daemon = map[string]string{}
	}
	c.Daemon[key] = value
	return write(root, c)
}

func write(root string, c Config) error {
	var b strings.Builder
	b.WriteString("# The Rhiza root: its projects and machine-wide settings.\n")
	if err := toml.NewEncoder(&b).Encode(c); err != nil {
		return err
	}
	path := filepath.Join(root, File)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
