package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/guide"
)

var pluginPaths = []string{pluginManifestPath, pluginSkillPath, pluginHooksPath, pluginBuildPath, pluginPlanPath, pluginScrubPath}

func guideInstall(t *testing.T, ctx *Ctx, req string) GuideInstallResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "guide.install", req)
	if err != nil {
		t.Fatal(err)
	}
	return resp.(GuideInstallResponse)
}

func writePath(t *testing.T, p, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestYAMLQuote(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"plain", `"plain"`},
		{"a: b", `"a: b"`},
		{`say "hi" \ bye`, `"say \"hi\" \\ bye"`},
		{"two\nlines\ttab", `"two\nlines\ttab"`},
		{"bell\x07", `"bell\x07"`},
		{"- dash # hash", `"- dash # hash"`},
	} {
		if got := yamlQuote(tc.in); got != tc.want {
			t.Errorf("yamlQuote(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

func TestSkillFileFrontmatter(t *testing.T) {
	for _, tc := range []struct{ name, text string }{{"rhiza", "rhiza"}, {"rhiza-scrub", "scrub"}, {"rhiza-plan", "plan"}, {"rhiza-build", "build"}} {
		b, err := skillFile(tc.name, tc.text)
		if err != nil {
			t.Fatal(err)
		}
		g, _ := guide.Lookup(tc.text)
		want := "---\nname: " + tc.name + "\ndescription: " + yamlQuote(g.Description) + "\n---\n\n" + g.Body
		if string(b) != want {
			t.Errorf("skillFile(%s) differs from the frontmatter plus the text", tc.name)
		}
	}
}

func TestGuideInstall(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	ctx.Version = "v1.2.3"
	dir := filepath.Join(ctx.Home, ".claude", "skills", "rhiza")

	r := guideInstall(t, ctx, `{"agent":"claude"}`)
	if r.Dir != dir || !slices.Equal(r.Written, pluginPaths) || len(r.Unchanged) != 0 {
		t.Fatalf("first install = %+v", r)
	}
	if v, ok := installedVersion(dir); !ok || v != "v1.2.3" {
		t.Errorf("installedVersion = %q, %v", v, ok)
	}

	r = guideInstall(t, ctx, `{"agent":"claude"}`)
	if len(r.Written) != 0 || !slices.Equal(r.Unchanged, pluginPaths) {
		t.Fatalf("second install = %+v", r)
	}

	ctx.Version = "v1.2.4"
	r = guideInstall(t, ctx, `{"agent":"claude"}`)
	if !slices.Equal(r.Written, []string{pluginManifestPath}) {
		t.Fatalf("new version install = %+v", r)
	}
	if v, _ := installedVersion(dir); v != "v1.2.4" {
		t.Errorf("installedVersion after upgrade = %q", v)
	}
}

// TestGuideInstallBuild finds the build topic in the installed plugin, as
// the rhiza-build skill: frontmatter, then rhiza guide build's text.
func TestGuideInstallBuild(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	r := guideInstall(t, ctx, `{"agent":"claude"}`)
	if !slices.Contains(r.Written, "skills/rhiza-build/SKILL.md") {
		t.Fatalf("install wrote %v, want skills/rhiza-build/SKILL.md", r.Written)
	}
	b, err := os.ReadFile(filepath.Join(r.Dir, "skills", "rhiza-build", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	g, _ := guide.Lookup("build")
	if !strings.HasPrefix(string(b), "---\nname: rhiza-build\ndescription: ") || !strings.HasSuffix(string(b), "\n---\n\n"+g.Body) {
		t.Errorf("rhiza-build SKILL.md is not the frontmatter plus the build topic:\n%s", b)
	}
	for _, p := range []string{"Leaf: ", "discovered-from", "rhiza run findings --leaf"} {
		if !strings.Contains(string(b), p) {
			t.Errorf("rhiza-build SKILL.md: no %q", p)
		}
	}
}

func TestGuideInstallForeign(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	dir := filepath.Join(ctx.Home, ".claude", "skills", "rhiza")
	writePath(t, filepath.Join(dir, "SKILL.md"), "mine\n")
	writePath(t, filepath.Join(dir, "extra.txt"), "keep\n")

	_, err := callOp(t, ctx, "guide.install", `{"agent":"claude"}`)
	var e *Error
	if !errors.As(err, &e) || e.Code != CodeConflict || e.Next != "rhiza guide install --agent claude --force" {
		t.Fatalf("foreign install err = %v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "SKILL.md")); string(b) != "mine\n" {
		t.Errorf("SKILL.md touched: %q", b)
	}
	if _, ok := installedVersion(dir); ok {
		t.Error("installedVersion ok on a foreign folder")
	}

	r := guideInstall(t, ctx, `{"agent":"claude","force":true}`)
	if !slices.Equal(r.Written, pluginPaths) {
		t.Fatalf("forced install = %+v", r)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "extra.txt")); string(b) != "keep\n" {
		t.Errorf("extra.txt = %q", b)
	}
}

func TestInstalledVersion(t *testing.T) {
	for _, tc := range []struct {
		name, manifest, version string
		ok                      bool
	}{
		{"Ours", `{"name":"rhiza","version":"v1","rhiza":{"installed_by":"rhiza guide install"}}`, "v1", true},
		{"EmptyMarker", `{"name":"rhiza","version":"v1","rhiza":{}}`, "v1", true},
		{"NoMarker", `{"name":"rhiza","version":"v1"}`, "", false},
		{"MarkerNotObject", `{"name":"rhiza","version":"v1","rhiza":"yes"}`, "", false},
		{"OtherName", `{"name":"other","version":"v1","rhiza":{}}`, "", false},
		{"NotJSON", `{`, "", false},
		{"Missing", "", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.manifest != "" {
				writePath(t, filepath.Join(dir, filepath.FromSlash(pluginManifestPath)), tc.manifest)
			}
			v, ok := installedVersion(dir)
			if v != tc.version || ok != tc.ok {
				t.Errorf("installedVersion = %q, %v; want %q, %v", v, ok, tc.version, tc.ok)
			}
		})
	}
}

func TestGuideInstallErrors(t *testing.T) {
	ctx := testCtx(t.TempDir())
	for _, tc := range []struct {
		name, req string
		dryRun    bool
		code      Code
		next      string
	}{
		{"OtherAgent", `{"agent":"agents-md"}`, false, CodeUsage, "rhiza guide install --agent claude"},
		{"NoHome", `{"agent":"claude"}`, false, CodeUsage, "rhiza guide install --agent claude --project"},
		{"ProjectNoRepo", `{"agent":"claude","project":true}`, false, CodeNotFound, "rhiza init"},
		{"DryRun", `{"agent":"claude","project":true}`, true, CodeUsage, "rhiza guide install --agent claude --project"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx.DryRun = tc.dryRun
			_, err := callOp(t, ctx, "guide.install", tc.req)
			var e *Error
			if !errors.As(err, &e) || e.Code != tc.code || e.Next != tc.next {
				t.Fatalf("err = %v", err)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(ctx.Root, ".claude")); err == nil {
		t.Error("an error wrote .claude")
	}
	if strings.Contains(pluginDescription, "\n") {
		t.Error("plugin description has a newline")
	}
}

func TestGuideInstallEmptyFolder(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	dir := filepath.Join(ctx.Home, ".claude", "skills", "rhiza")
	writePath(t, filepath.Join(dir, "notes.txt"), "keep\n")

	r := guideInstall(t, ctx, `{"agent":"claude"}`)
	if !slices.Equal(r.Written, pluginPaths) {
		t.Errorf("written = %v, want %v", r.Written, pluginPaths)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "notes.txt")); string(b) != "keep\n" {
		t.Errorf("notes.txt touched: %q", b)
	}
}

func TestGuideInstallFileInTheWay(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	dir := filepath.Join(ctx.Home, ".claude", "skills", "rhiza")
	writePath(t, dir, "a file\n")

	for _, req := range []string{`{"agent":"claude"}`, `{"agent":"claude","force":true}`} {
		_, err := callOp(t, ctx, "guide.install", req)
		var e *Error
		if !errors.As(err, &e) || e.Code != CodeConflict || e.Next != "" || !strings.Contains(e.Message, "is a file, not a folder") {
			t.Errorf("%s: err = %v", req, err)
		}
	}
}

// TestGuideInstallConcurrent runs first installs at once: each sees the
// folder another may have just created, and none may call it foreign.
func TestGuideInstallConcurrent(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Home = t.TempDir()
	const n = 8
	errs := make(chan error, n)
	o, _ := Lookup("guide.install")
	req, err := DecodeRequest(o, []byte(`{"agent":"claude"}`))
	if err != nil {
		t.Fatal(err)
	}
	for range n {
		go func() {
			_, err := o.Run(ctx, req)
			errs <- err
		}()
	}
	for range n {
		if err := <-errs; err != nil {
			t.Error(err)
		}
	}
	if v, ok := installedVersion(filepath.Join(ctx.Home, ".claude", "skills", "rhiza")); !ok || v != ctx.Version {
		t.Errorf("installedVersion = %q, %v", v, ok)
	}
}
