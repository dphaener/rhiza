// Package anchor resolves a code anchor against the working tree: a path
// relative to the project root and an optional Go symbol, and fingerprints
// what they name.
package anchor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Prefix starts every fingerprint.
const Prefix = "sha256:"

// The ways an anchor fails to resolve.
var (
	// ErrBadPath is a path that is not relative, slash-separated and inside
	// the project root.
	ErrBadPath = errors.New("bad path")
	// ErrNoPath is a path with no regular file in the working tree.
	ErrNoPath = errors.New("path not found")
	// ErrBadSymbol is a symbol on a file that is not Go source, or a Go file
	// that does not parse.
	ErrBadSymbol = errors.New("bad symbol")
	// ErrNoSymbol is a symbol the Go file does not declare at top level.
	ErrNoSymbol = errors.New("symbol not found")
)

// Fingerprint returns the fingerprint of what path and symbol name under
// root: sha256 of the declaration's source bytes when symbol is set, else of
// the whole file's bytes. The error wraps one of the Err values.
func Fingerprint(root, p, symbol string) (string, error) {
	if err := checkPath(p); err != nil {
		return "", err
	}
	full := filepath.Join(root, filepath.FromSlash(p))
	info, err := os.Stat(full)
	if err != nil || !info.Mode().IsRegular() {
		if err == nil || errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrInvalid) {
			return "", fmt.Errorf("%w: %s", ErrNoPath, p)
		}
		return "", fmt.Errorf("%w: %s: %v", ErrNoPath, p, err)
	}
	src, err := os.ReadFile(full)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %v", ErrNoPath, p, err)
	}
	if symbol == "" {
		return hash(src), nil
	}
	if !strings.HasSuffix(p, ".go") {
		return "", fmt.Errorf("%w: %s is not a .go file, so it has no symbols", ErrBadSymbol, p)
	}
	decl, err := declaration(src, symbol)
	if err != nil {
		return "", fmt.Errorf("%s: %w", p, err)
	}
	return hash(decl), nil
}

func hash(b []byte) string {
	sum := sha256.Sum256(b)
	return Prefix + hex.EncodeToString(sum[:])
}

func checkPath(p string) error {
	switch {
	case p == "":
		return fmt.Errorf("%w: empty", ErrBadPath)
	case strings.Contains(p, `\`):
		return fmt.Errorf("%w: %q has a backslash; use slashes", ErrBadPath, p)
	case path.IsAbs(p):
		return fmt.Errorf("%w: %q is absolute; use a path relative to the project root", ErrBadPath, p)
	case path.Clean(p) != p || p == "." || p == ".." || strings.HasPrefix(p, "../"):
		return fmt.Errorf("%w: %q is not a clean path inside the project root", ErrBadPath, p)
	}
	return nil
}

// declaration returns the source bytes of the top-level declaration symbol
// names. A func or a type, var or const outside a group is the whole
// declaration; one inside a group is its own spec. A doc comment is not part
// of it.
func declaration(src []byte, symbol string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "anchor.go", src, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("%w: does not parse: %v", ErrBadSymbol, err)
	}
	span := func(n ast.Node) []byte {
		return src[fset.Position(n.Pos()).Offset:fset.Position(n.End()).Offset]
	}
	recv, name, isMethod := strings.Cut(symbol, ".")
	if !isMethod {
		recv, name = "", symbol
	}
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Name.Name != name {
				continue
			}
			if isMethod != (d.Recv != nil) || (isMethod && receiverName(d.Recv) != recv) {
				continue
			}
			return span(d), nil
		case *ast.GenDecl:
			if isMethod {
				continue
			}
			for _, s := range d.Specs {
				if !declares(s, name) {
					continue
				}
				if d.Lparen.IsValid() {
					return span(s), nil
				}
				return span(d), nil
			}
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrNoSymbol, symbol)
}

func declares(s ast.Spec, name string) bool {
	switch s := s.(type) {
	case *ast.TypeSpec:
		return s.Name.Name == name
	case *ast.ValueSpec:
		for _, n := range s.Names {
			if n.Name == name {
				return true
			}
		}
	}
	return false
}

// receiverName is the type name of a method's receiver, without the pointer
// star or type parameters.
func receiverName(r *ast.FieldList) string {
	if r == nil || len(r.List) == 0 {
		return ""
	}
	t := r.List[0].Type
	for {
		switch x := t.(type) {
		case *ast.StarExpr:
			t = x.X
		case *ast.ParenExpr:
			t = x.X
		case *ast.IndexExpr:
			t = x.X
		case *ast.IndexListExpr:
			t = x.X
		case *ast.Ident:
			return x.Name
		default:
			return ""
		}
	}
}
