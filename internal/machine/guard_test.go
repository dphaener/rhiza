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

// importNames maps each name a file uses for an import to its path, so an
// aliased import ("xos" for "os") is caught like the plain one.
func importNames(f *ast.File) map[string]string {
	names := map[string]string{}
	for _, im := range f.Imports {
		p, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			continue
		}
		name := p[strings.LastIndex(p, "/")+1:]
		if im.Name != nil {
			name = im.Name.Name
		}
		names[name] = p
	}
	return names
}

// pkgCall returns the import path and function name of a call to a
// package-level function, or "" when the call is not one.
func pkgCall(imports map[string]string, call *ast.CallExpr) (string, string) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", ""
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Obj != nil {
		return "", ""
	}
	return imports[id.Name], sel.Sel.Name
}

// readOnly reports whether an os.OpenFile flag argument is exactly
// os.O_RDONLY: any other flag, or one held in a variable, may write.
func readOnly(imports map[string]string, flag ast.Expr) bool {
	sel, ok := flag.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && imports[id.Name] == "os" && sel.Sel.Name == "O_RDONLY"
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
		imports := importNames(f)
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			pkg, fn := pkgCall(imports, call)
			switch {
			case (pkg == "os" || pkg == "io/ioutil") && writers[fn]:
				t.Errorf("%s calls %s.%s: internal/machine may only read", name, pkg, fn)
			case pkg == "os" && fn == "OpenFile" && (len(call.Args) < 2 || !readOnly(imports, call.Args[1])):
				t.Errorf("%s calls os.OpenFile without os.O_RDONLY: internal/machine may only read", name)
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
	fileAPIs := map[string]bool{"path/filepath": true, "path": true, "os": true, "io/ioutil": true}
	for name, f := range files {
		imports := importNames(f)
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			pkg, fn := pkgCall(imports, call)
			if !fileAPIs[pkg] {
				return true
			}
			for _, a := range call.Args {
				ast.Inspect(a, func(m ast.Node) bool {
					if l, ok := m.(*ast.BasicLit); ok && l.Kind == token.STRING {
						if s, err := strconv.Unquote(l.Value); err == nil && strings.Contains(s, "machine.toml") {
							t.Errorf("%s passes %q to %s.%s: the machine config is reached only through internal/machine", name, s, pkg, fn)
						}
					}
					return true
				})
			}
			return true
		})
	}
}
