package doclinks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, root, name, body string) {
	t.Helper()
	p := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBrokenLinksFail(t *testing.T) {
	root := t.TempDir()
	write(t, root, "README.md", "# Top\n\n[ok](docs/a.md#the-install-step) [file](docs/gone.md) [anchor](docs/a.md#nope) [self](#top) [bad self](#missing) [web](https://example.com/x.md)\n")
	write(t, root, "docs/a.md", "# A\n\n## The `install` step\n\n[up](../README.md) [sib](b.md)\n\n```bash\n[ignored](nothing.md)\n```\n")
	got, err := Check(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"README.md:3: link \"#missing\": no heading #missing in README.md",
		"README.md:3: link \"docs/a.md#nope\": no heading #nope in docs/a.md",
		"README.md:3: link \"docs/gone.md\": docs/gone.md does not exist",
		"docs/a.md:5: link \"b.md\": docs/b.md does not exist",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("got:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestCleanFixturePasses(t *testing.T) {
	root := t.TempDir()
	write(t, root, "README.md", "# Top\n\n[q](docs/q.md#where-next)\n")
	write(t, root, "docs/q.md", "# Q\n\n## Where next\n\n[home](../README.md#top)\n")
	got, err := Check(root)
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v", got, err)
	}
}

// The real README and docs/ must have no broken relative link. This runs in
// make check through go test.
func TestRepoDocsLinksResolve(t *testing.T) {
	got, err := Check(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("broken doc links:\n%s", strings.Join(got, "\n"))
	}
}
