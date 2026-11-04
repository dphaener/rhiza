package anchor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const src = `package p

// Doc for F.
func F() int { return 1 }

type T struct{ x int }

// Doc for V.
func (t T) V() int { return t.x }

func (t *T) P() { t.x++ }

type G[K any] struct{}

func (g *G[K]) M() {}

var Solo = 1

var (
	A = 1
	B, C = 2, 3
)

const (
	K1 = "a"
	K2 = "b"
)

type (
	Alpha int
	Beta  string
)
`

func sum(s string) string {
	h := sha256.Sum256([]byte(s))
	return Prefix + hex.EncodeToString(h[:])
}

func TestFingerprintSymbols(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "p.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"F":    "func F() int { return 1 }",
		"T":    "type T struct{ x int }",
		"T.V":  "func (t T) V() int { return t.x }",
		"T.P":  "func (t *T) P() { t.x++ }",
		"G.M":  "func (g *G[K]) M() {}",
		"G":    "type G[K any] struct{}",
		"Solo": "var Solo = 1",
		"A":    "A = 1",
		"C":    "B, C = 2, 3",
		"K2":   `K2 = "b"`,
		"Beta": "Beta  string",
	}
	for sym, want := range cases {
		got, err := Fingerprint(root, "p.go", sym)
		if err != nil {
			t.Errorf("%s: %v", sym, err)
			continue
		}
		if got != sum(want) {
			t.Errorf("%s: got %s, want hash of %q", sym, got, want)
		}
	}
}

func TestFingerprintFile(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "docs"), 0o755)
	os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("hello\n"), 0o644)
	got, err := Fingerprint(root, "docs/a.md", "")
	if err != nil || got != sum("hello\n") {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestFingerprintErrors(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "p.go"), []byte(src), 0o644)
	os.WriteFile(filepath.Join(root, "a.md"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(root, "bad.go"), []byte("package ("), 0o644)
	os.Mkdir(filepath.Join(root, "dir"), 0o755)
	cases := []struct {
		path, sym string
		want      error
	}{
		{"nope.go", "", ErrNoPath},
		{"nope.go", "F", ErrNoPath},
		{"dir", "", ErrNoPath},
		{"p.go", "Missing", ErrNoSymbol},
		{"p.go", "T.Missing", ErrNoSymbol},
		{"p.go", "F.V", ErrNoSymbol},
		{"p.go", "V", ErrNoSymbol},
		{"a.md", "F", ErrBadSymbol},
		{"bad.go", "F", ErrBadSymbol},
		{"../p.go", "", ErrBadPath},
		{"/etc/passwd", "", ErrBadPath},
		{"", "", ErrBadPath},
		{`a\b`, "", ErrBadPath},
	}
	for _, c := range cases {
		_, err := Fingerprint(root, c.path, c.sym)
		if !errors.Is(err, c.want) {
			t.Errorf("%q %q: got %v, want %v", c.path, c.sym, err, c.want)
		}
	}
}
