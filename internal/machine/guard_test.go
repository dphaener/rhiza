package machine

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// The machine config is changed only by hand: rhiza must have no code path
// that writes or raises a machine's trust. These guards parse the module's
// non-test source and fail when this package could write the file, widens
// its API, or when another package reaches the file itself.

// writers are the os functions that create, write or change files.
var writers = map[string]bool{
	"WriteFile": true, "Create": true, "CreateTemp": true, "Mkdir": true,
	"MkdirAll": true, "MkdirTemp": true, "Rename": true, "Remove": true,
	"RemoveAll": true, "Chmod": true, "Chown": true, "Chtimes": true,
	"Symlink": true, "Link": true, "Truncate": true,
}

// writeFlags are the os.OpenFile flags that make a file writable.
var writeFlags = map[string]bool{
	"O_WRONLY": true, "O_RDWR": true, "O_APPEND": true, "O_CREATE": true,
	"O_TRUNC": true,
}

func parseDir(t *testing.T, root string, skip func(rel string) bool) map[string]*ast.File {
	t.Helper()
	fset := token.NewFileSet()
	files := map[string]*ast.File{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if d.IsDir() {
			if rel != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata" || skip(rel)) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = f
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestMachineNeverWritesFiles(t *testing.T) {
	files := parseDir(t, ".", func(string) bool { return true })
	if len(files) == 0 {
		t.Fatal("no source parsed")
	}
	for name, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != "os" {
				return true
			}
			switch {
			case writers[sel.Sel.Name]:
				t.Errorf("%s calls os.%s: internal/machine may only read", name, sel.Sel.Name)
			case sel.Sel.Name == "OpenFile":
				ast.Inspect(call, func(m ast.Node) bool {
					if s, ok := m.(*ast.SelectorExpr); ok {
						if p, ok := s.X.(*ast.Ident); ok && p.Name == "os" && writeFlags[s.Sel.Name] {
							t.Errorf("%s opens a file with os.%s", name, s.Sel.Name)
						}
					}
					return true
				})
			}
			return true
		})
	}
}

func TestMachineExportsOnlyReaders(t *testing.T) {
	want := []string{"Config", "Daily", "Dedicated", "Load", "Path", "Read"}
	var got []string
	for name, f := range parseDir(t, ".", func(string) bool { return true }) {
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if d.Name.IsExported() {
					id := d.Name.Name
					if d.Recv != nil {
						id = "method " + id
					}
					got = append(got, id)
				}
			case *ast.GenDecl:
				for _, s := range d.Specs {
					switch s := s.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							got = append(got, s.Name.Name)
						}
					case *ast.ValueSpec:
						for _, id := range s.Names {
							if id.IsExported() {
								got = append(got, id.Name)
							}
						}
					}
				}
			}
		}
		_ = name
	}
	sort.Strings(got)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("internal/machine exports %v, want %v", got, want)
	}
}

func TestNothingElseTouchesMachineFile(t *testing.T) {
	root := filepath.Join("..", "..")
	files := parseDir(t, root, func(rel string) bool {
		return filepath.ToSlash(rel) == "internal/machine"
	})
	if len(files) < 10 {
		t.Fatalf("parsed only %d files from %s", len(files), root)
	}
	fileAPIs := map[string]bool{"filepath": true, "path": true, "os": true, "ioutil": true}
	for name, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || !fileAPIs[pkg.Name] {
				return true
			}
			for _, a := range call.Args {
				ast.Inspect(a, func(m ast.Node) bool {
					if l, ok := m.(*ast.BasicLit); ok && l.Kind == token.STRING {
						if s, err := strconv.Unquote(l.Value); err == nil && strings.Contains(s, "machine.toml") {
							t.Errorf("%s passes %q to %s.%s: the machine config is reached only through internal/machine", name, s, pkg.Name, sel.Sel.Name)
						}
					}
					return true
				})
			}
			return true
		})
	}
}
