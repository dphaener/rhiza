// Package key provides fractional sibling keys: a key orders siblings, and
// inserting between two keys never renumbers the others.
//
// A key is a string of base-36 digits (0-9 then a-z), read as the digits of a
// fraction after the point. Keys compare as plain byte strings
// (strings.Compare), and that is the sibling order. A key never ends in "0",
// so there is always room for another key below it, and is at most MaxLen
// characters long.
//
// Between picks the digit-wise midpoint of two keys, as the fractional-indexing
// algorithm (Evan Wallace, after David Greenspan) does, without its
// integer-part header. It is deterministic: the same bounds give the same key.
package key

import (
	"errors"
	"fmt"
	"strings"
)

// MaxLen is the longest valid key.
const MaxLen = 256

// digits is the key alphabet, in ascending byte order.
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	// ErrInvalid is wrapped by every error about a malformed key.
	ErrInvalid = errors.New("invalid key")
	// ErrOrder is returned by Between when the lower bound is not below the
	// upper bound.
	ErrOrder = errors.New("keys out of order")
	// ErrTooLong is returned by Between when the key between its bounds would
	// be longer than MaxLen.
	ErrTooLong = errors.New("key too long")
)

// Valid reports whether k is a valid key: non-empty, at most MaxLen base-36
// digits, not ending in "0".
func Valid(k string) error {
	switch {
	case k == "":
		return fmt.Errorf("%w: empty", ErrInvalid)
	case len(k) > MaxLen:
		return fmt.Errorf("%w %.16q...: %d characters, more than %d", ErrInvalid, k, len(k), MaxLen)
	}
	for i := 0; i < len(k); i++ {
		if strings.IndexByte(digits, k[i]) < 0 {
			return fmt.Errorf("%w %q: character %q at %d is not a base-36 digit", ErrInvalid, k, k[i], i)
		}
	}
	if k[len(k)-1] == '0' {
		return fmt.Errorf("%w %q: ends in 0", ErrInvalid, k)
	}
	return nil
}

// Between returns a valid key k with a < k < b. An empty a means no lower
// bound and an empty b no upper bound; a bound that is given must be valid,
// and a must be below b when both are given. The result is the digit-wise
// midpoint, so it is as short as that gives; Between("", "") is "i".
func Between(a, b string) (string, error) {
	if a != "" {
		if err := Valid(a); err != nil {
			return "", fmt.Errorf("lower bound: %w", err)
		}
	}
	if b != "" {
		if err := Valid(b); err != nil {
			return "", fmt.Errorf("upper bound: %w", err)
		}
		if a >= b {
			return "", fmt.Errorf("%w: %q is not below %q", ErrOrder, a, b)
		}
	}
	k := midpoint(a, b)
	if len(k) > MaxLen {
		return "", fmt.Errorf("%w: between %q and %q needs %d characters, more than %d", ErrTooLong, a, b, len(k), MaxLen)
	}
	return k, nil
}

// midpoint returns the key halfway between a and b. a may be empty (zero);
// b empty means one. a < b, and neither ends in "0".
func midpoint(a, b string) string {
	var out strings.Builder
	for {
		if b != "" {
			// Copy the common prefix, reading a as padded with zeros.
			n := 0
			for n < len(b) && digitAt(a, n) == b[n] {
				n++
			}
			out.WriteString(b[:n])
			a, b = a[min(n, len(a)):], b[n:]
		}
		da := 0
		if a != "" {
			da = strings.IndexByte(digits, a[0])
		}
		db := len(digits)
		if b != "" {
			db = strings.IndexByte(digits, b[0])
		}
		if db-da > 1 {
			out.WriteByte(digits[(da+db+1)/2])
			return out.String()
		}
		// The first digits are consecutive.
		if len(b) > 1 {
			// b's first digit alone is above a and below b.
			out.WriteByte(b[0])
			return out.String()
		}
		// Keep a's first digit and find a key above the rest of a.
		out.WriteByte(digits[da])
		if a != "" {
			a = a[1:]
		}
		b = ""
	}
}

// digitAt returns k[i], or '0' past the end of k.
func digitAt(k string, i int) byte {
	if i < len(k) {
		return k[i]
	}
	return '0'
}
