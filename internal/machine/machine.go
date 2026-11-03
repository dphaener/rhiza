// Package machine reads the per-machine config: the one file that says how
// far this machine is trusted. It lives outside every repo, so no clone
// carries it.
package machine

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Trust levels a machine config may declare.
const (
	Daily     = "daily"
	Dedicated = "dedicated"
)

// Config is what the machine config file holds.
type Config struct {
	// Trust is Daily or Dedicated.
	Trust string `toml:"trust"`
}

// Path returns the machine config file: $XDG_CONFIG_HOME/rhiza/machine.toml,
// else <home>/.config/rhiza/machine.toml. It is empty when neither is known.
func Path(getenv func(string) string, home string) string {
	if getenv != nil {
		if x := getenv("XDG_CONFIG_HOME"); x != "" {
			return filepath.Join(x, "rhiza", "machine.toml")
		}
	}
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "rhiza", "machine.toml")
}

// Read reads the machine config at path. An absent file (or an empty path)
// is trust daily; a file that is not TOML or names an unknown trust is an
// error naming path. Nothing is cached: call it at each check.
func Read(path string) (Config, error) {
	if path == "" {
		return Config{Trust: Daily}, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{Trust: Daily}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("machine config %s: %w", path, err)
	}
	var c Config
	if _, err := toml.Decode(string(data), &c); err != nil {
		return Config{}, fmt.Errorf("machine config %s: %w", path, err)
	}
	switch c.Trust {
	case "":
		c.Trust = Daily
	case Daily, Dedicated:
	default:
		return Config{}, fmt.Errorf("machine config %s: trust = %q: want \"dedicated\" or \"daily\"", path, c.Trust)
	}
	return c, nil
}

// Load finds the machine config from the environment and reads it.
func Load(getenv func(string) string, home string) (Config, error) {
	return Read(Path(getenv, home))
}
