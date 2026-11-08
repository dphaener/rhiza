// Package policy reads a repo's scheduling rule from its checked-in
// .rhiza/policy.toml:
//
//	[schedule]
//	backlog = "after"
//
// backlog says what the daemon does with accepted change sets that serve no
// milestone: "after" builds them after milestone work (the default, also
// used when the file or the key is missing), "before" builds them before it
// and "skip" never builds them. Unknown keys are ignored.
package policy

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

// Path is the policy file, relative to the project root.
const Path = ".rhiza/policy.toml"

// The values of the schedule's backlog key.
const (
	BacklogAfter  = "after"
	BacklogBefore = "before"
	BacklogSkip   = "skip"
)

// Policy is a repo's policy.
type Policy struct {
	// Backlog is one of BacklogAfter, BacklogBefore and BacklogSkip, never
	// empty.
	Backlog string
}

// Default is the policy of a repo without a policy file.
func Default() Policy { return Policy{Backlog: BacklogAfter} }

// file is the shape of policy.toml. Pointers tell a missing key from "".
type file struct {
	Schedule struct {
		Backlog *string `toml:"backlog"`
	} `toml:"schedule"`
}

// Read reads the policy from fsys, rooted at the project. A missing file
// or key is the default. A file that does not parse, or a backlog value
// outside the three, is an error naming Path (and the value); the returned
// Policy is then the default.
func Read(fsys fs.FS) (Policy, error) {
	p := Default()
	b, err := fs.ReadFile(fsys, Path)
	if errors.Is(err, fs.ErrNotExist) {
		return p, nil
	}
	if err != nil {
		return p, fmt.Errorf("%s: %w", Path, err)
	}
	var f file
	if _, err := toml.Decode(string(b), &f); err != nil {
		return p, fmt.Errorf("%s: %w", Path, err)
	}
	if v := f.Schedule.Backlog; v != nil {
		switch *v {
		case BacklogAfter, BacklogBefore, BacklogSkip:
			p.Backlog = *v
		default:
			return p, fmt.Errorf("%s: schedule.backlog %q is not one of %q, %q or %q",
				Path, *v, BacklogAfter, BacklogBefore, BacklogSkip)
		}
	}
	return p, nil
}
