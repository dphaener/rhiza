// Package id provides Genesis's identifiers: format-1 ULIDs, random and
// time-sortable, generated anywhere without coordination.
//
// Every ID has exactly one spelling, its canonical 26-character upper-case
// Crockford base32 form. Log entries are canonical bytes, so Parse rejects
// every other spelling (lower case included) rather than normalising it.
package id

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// ErrInvalid is wrapped by every error Parse and UnmarshalText return.
var ErrInvalid = errors.New("invalid id")

// ID is a ULID. The zero value is the all-zero ULID and is never generated.
type ID ulid.ULID

// String returns the canonical 26-character upper-case Crockford base32 form.
func (i ID) String() string {
	return ulid.ULID(i).String()
}

// IsZero reports whether i is the zero ID.
func (i ID) IsZero() bool {
	return i == ID{}
}

// Parse returns the ID spelled by s, which must be in canonical form.
func Parse(s string) (ID, error) {
	u, err := ulid.ParseStrict(s)
	if err != nil {
		return ID{}, fmt.Errorf("%w %q: %w", ErrInvalid, s, err)
	}
	// ParseStrict accepts lower case; one spelling per ID means we don't.
	if u.String() != s {
		return ID{}, fmt.Errorf("%w %q: not canonical upper case", ErrInvalid, s)
	}
	return ID(u), nil
}

// Compare returns -1, 0 or +1 as a sorts before, equal to or after b. Byte
// order and canonical string order agree.
func Compare(a, b ID) int {
	return bytes.Compare(a[:], b[:])
}

// MarshalText implements encoding.TextMarshaler with the canonical form.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler with Parse's rules.
func (i *ID) UnmarshalText(b []byte) error {
	parsed, err := Parse(string(b))
	if err != nil {
		return err
	}
	*i = parsed
	return nil
}

// Generator makes IDs that strictly increase within one millisecond. It is
// safe for concurrent use.
type Generator struct {
	mu      sync.Mutex
	now     func() time.Time
	entropy *ulid.MonotonicEntropy
}

// NewGenerator returns a Generator reading the time from now and randomness
// from entropy. Tests pass a fixed clock and seeded entropy for repeatable IDs.
func NewGenerator(now func() time.Time, entropy io.Reader) *Generator {
	return &Generator{now: now, entropy: ulid.Monotonic(entropy, 0)}
}

// New returns the next ID. It fails if entropy fails or the monotonic
// sequence for the current millisecond is exhausted.
func (g *Generator) New() (ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	u, err := ulid.New(ulid.Timestamp(g.now()), g.entropy)
	if err != nil {
		return ID{}, fmt.Errorf("generate id: %w", err)
	}
	return ID(u), nil
}

var defaultGenerator = NewGenerator(time.Now, rand.Reader)

// Default returns the process-wide Generator, using the system clock and
// crypto/rand.
func Default() *Generator {
	return defaultGenerator
}
