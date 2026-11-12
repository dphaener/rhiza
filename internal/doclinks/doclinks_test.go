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

// The daemon docs name the autoUpdate setting and its values, and no longer
// the removed selfUpdate setting or self_update_project key.
func TestDocsDescribeAutoUpdate(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "docs", "runs-and-daemon.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, want := range []string{"rhiza.daemon.autoUpdate", "`off`", "`main`"} {
		if !strings.Contains(doc, want) {
			t.Errorf("docs/runs-and-daemon.md does not name %s", want)
		}
	}
	for _, gone := range []string{"selfUpdate", "self_update_project"} {
		if strings.Contains(doc, gone) {
			t.Errorf("docs/runs-and-daemon.md still names %s", gone)
		}
	}
}
