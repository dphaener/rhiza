package key

import (
	"errors"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		k    string
		ok   bool
	}{
		{"single digit", "i", true},
		{"lowest one-digit key", "1", true},
		{"highest one-digit key", "z", true},
		{"inner zeros", "a00b", true},
		{"all digits", "0123456789abcdefghijklmnopqrstuvwxyz", true},
		{"max length", strings.Repeat("z", MaxLen), true},
		{"empty", "", false},
		{"bad character", "a-b", false},
		{"space", "a b", false},
		{"non-ascii", "aé", false},
		{"upper case", "A", false},
		{"upper case inside", "aBc", false},
		{"trailing 0", "a0", false},
		{"only 0", "0", false},
		{"too long", strings.Repeat("z", MaxLen+1), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Valid(tc.k)
			if tc.ok {
				if err != nil {
					t.Fatalf("Valid(%q) = %v, want nil", tc.k, err)
				}
				return
			}
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Valid(%q) = %v, want ErrInvalid", tc.k, err)
			}
			if tc.k != "" && len(tc.k) < 40 && !strings.Contains(err.Error(), tc.k) {
				t.Fatalf("Valid(%q) error %q does not name the key", tc.k, err)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want string
	}{
		{"both bounds open", "", "", "i"},
		{"lower bound open", "", "i", "9"},
		{"lower bound open below 1", "", "1", "0i"},
		{"lower bound open below 01", "", "01", "00i"},
		{"upper bound open", "i", "", "r"},
		{"upper bound open above z", "z", "", "zi"},
		{"upper bound open above zz", "zz", "", "zzi"},
		{"upper bound open above y", "y", "", "z"},
		{"wide gap", "a", "z", "n"},
		{"adjacent keys", "a", "b", "ai"},
		{"adjacent keys at the top", "y", "z", "yi"},
		{"below z", "yz", "z", "yzi"},
		{"different lengths: a, a1", "a", "a1", "a0i"},
		{"different lengths: a, a01", "a", "a01", "a00i"},
		{"different lengths: az, b", "az", "b", "azi"},
		{"different lengths: a, bc", "a", "bc", "b"},
		{"different lengths: ab, b", "ab", "b", "ao"},
		{"common prefix", "abc", "abe", "abd"},
		{"common prefix adjacent", "abc", "abd", "abci"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Between(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Between(%q, %q) = %v", tc.a, tc.b, err)
			}
			if got != tc.want {
				t.Fatalf("Between(%q, %q) = %q, want %q", tc.a, tc.b, got, tc.want)
			}
			checkBetween(t, tc.a, tc.b, got)
		})
	}
}

func TestBetweenErrors(t *testing.T) {
	long := strings.Repeat("z", MaxLen)
	tests := []struct {
		name string
		a, b string
		want error
	}{
		{"upper case lower bound", "A", "b", ErrInvalid},
		{"upper case upper bound", "a", "B", ErrInvalid},
		{"bad character lower bound", "a.", "", ErrInvalid},
		{"bad character upper bound", "", "-", ErrInvalid},
		{"trailing 0 lower bound", "a0", "b", ErrInvalid},
		{"trailing 0 upper bound", "", "b0", ErrInvalid},
		{"too long lower bound", strings.Repeat("a", MaxLen+1), "", ErrInvalid},
		{"too long upper bound", "", strings.Repeat("a", MaxLen+1), ErrInvalid},
		{"a >= b equal", "a", "a", ErrOrder},
		{"a >= b reversed", "b", "a", ErrOrder},
		{"a >= b reversed prefix", "a1", "a", ErrOrder},
		{"too long result above", long, "", ErrTooLong},
		{"too long result between", long[:MaxLen-1] + "y", long, ErrTooLong},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Between(tc.a, tc.b)
			if !errors.Is(err, tc.want) {
				t.Fatalf("Between(%q, %q) = %q, %v; want %v", tc.a, tc.b, got, err, tc.want)
			}
			if got != "" {
				t.Fatalf("Between(%q, %q) returned key %q with error", tc.a, tc.b, got)
			}
		})
	}
}

// checkBetween fails unless k is a valid key strictly between the bounds.
func checkBetween(t interface {
	Helper()
	Fatalf(string, ...any)
}, a, b, k string) {
	t.Helper()
	if err := Valid(k); err != nil {
		t.Fatalf("Between(%q, %q) = %q: %v", a, b, k, err)
	}
	if a != "" && !(a < k) {
		t.Fatalf("Between(%q, %q) = %q, not above the lower bound", a, b, k)
	}
	if b != "" && !(k < b) {
		t.Fatalf("Between(%q, %q) = %q, not below the upper bound", a, b, k)
	}
}

// genKey draws a valid key of up to maxLen characters.
func genKey(maxLen int) *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		body := rapid.StringOfN(rapid.RuneFrom([]rune(digits)), 0, maxLen-1, -1).Draw(t, "body")
		last := rapid.RuneFrom([]rune(digits[1:])).Draw(t, "last")
		return body + string(last)
	})
}

func TestBetweenProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Mostly short keys with shared prefixes; sometimes max-length ones.
		maxLen := rapid.SampledFrom([]int{3, 8, MaxLen}).Draw(t, "maxLen")
		x := genKey(maxLen).Draw(t, "x")
		y := genKey(maxLen).Draw(t, "y")
		if x == y {
			y = ""
		}
		a, b := min(x, y), max(x, y)
		if y == "" {
			a, b = x, ""
		}
		switch rapid.IntRange(0, 3).Draw(t, "bounds") {
		case 1:
			a = ""
		case 2:
			b = ""
		case 3:
			a, b = "", ""
		}
		k, err := Between(a, b)
		if errors.Is(err, ErrTooLong) && max(len(a), len(b)) == MaxLen {
			return
		}
		if err != nil {
			t.Fatalf("Between(%q, %q) = %v", a, b, err)
		}
		checkBetween(t, a, b, k)
		if again, _ := Between(a, b); again != k {
			t.Fatalf("Between(%q, %q) not deterministic: %q then %q", a, b, k, again)
		}
	})
}

func TestRepeatedInsertProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		lo := genKey(6).Draw(t, "lo")
		hi := genKey(6).Draw(t, "hi")
		if lo == hi {
			hi = lo + "1"
		}
		lo, hi = min(lo, hi), max(lo, hi)
		mode := rapid.SampledFrom([]string{"front", "back", "between"}).Draw(t, "mode")
		// keys is the sibling list; every insert must keep it sorted.
		keys := []string{lo, hi}
		for i := range 1000 {
			var a, b string
			var at int
			switch mode {
			case "front":
				a, b, at = "", keys[0], 0
			case "back":
				a, b, at = keys[len(keys)-1], "", len(keys)
			case "between":
				// Always insert next to one end of the fixed pair, so the gap
				// keeps shrinking on that side.
				if rapid.Bool().Draw(t, "left") {
					a, b, at = keys[0], keys[1], 1
				} else {
					a, b, at = keys[len(keys)-2], keys[len(keys)-1], len(keys)-1
				}
			}
			k, err := Between(a, b)
			if errors.Is(err, ErrTooLong) {
				return
			}
			if err != nil {
				t.Fatalf("insert %d: Between(%q, %q) = %v", i, a, b, err)
			}
			checkBetween(t, a, b, k)
			keys = append(keys[:at], append([]string{k}, keys[at:]...)...)
		}
		for i := 1; i < len(keys); i++ {
			if !(keys[i-1] < keys[i]) {
				t.Fatalf("keys out of order at %d: %q, %q", i, keys[i-1], keys[i])
			}
		}
	})
}

func FuzzBetween(f *testing.F) {
	seeds := [][2]string{
		{"", ""}, {"", "i"}, {"i", ""}, {"a", "b"}, {"a", "a1"}, {"az", "b"},
		{"z", ""}, {"zz", ""}, {"", "1"}, {"a0", "b"}, {"A", ""}, {"b", "a"},
		{strings.Repeat("z", MaxLen), ""},
	}
	for _, s := range seeds {
		f.Add(s[0], s[1])
	}
	f.Fuzz(func(t *testing.T, a, b string) {
		k, err := Between(a, b)
		if err != nil {
			if k != "" {
				t.Fatalf("Between(%q, %q) returned %q with %v", a, b, k, err)
			}
			if !errors.Is(err, ErrInvalid) && !errors.Is(err, ErrOrder) && !errors.Is(err, ErrTooLong) {
				t.Fatalf("Between(%q, %q) = unexpected error %v", a, b, err)
			}
			goodBounds := (a == "" || Valid(a) == nil) && (b == "" || Valid(b) == nil)
			if errors.Is(err, ErrInvalid) == goodBounds {
				t.Fatalf("Between(%q, %q) = %v, bounds valid: %v", a, b, err, goodBounds)
			}
			if errors.Is(err, ErrOrder) && (a == "" || b == "" || a < b) {
				t.Fatalf("Between(%q, %q) = %v for ordered bounds", a, b, err)
			}
			if errors.Is(err, ErrTooLong) && max(len(a), len(b)) < MaxLen {
				t.Fatalf("Between(%q, %q) = %v for short bounds", a, b, err)
			}
			return
		}
		checkBetween(t, a, b, k)
	})
}
