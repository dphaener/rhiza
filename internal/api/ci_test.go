package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestInstallVersion(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"v1.2.3", "v1.2.3"},
		{"v0.3.0", "v0.3.0"},
		{"v1.2.3-rc.1", "v1.2.3-rc.1"},
		{"v0.0.0-20261024120000-abcdef123456", "v0.0.0-20261024120000-abcdef123456"},
		// git describe: commits past a tag name no published version.
		{"v1.2.3-4-gabcdef1", "latest"},
		{"v1.2.3-dirty", "latest"},
		{"v1.2.3-4-gabcdef1-dirty", "latest"},
		// A VCS-stamped module version with local edits.
		{"v0.0.0-20261024120000-abcdef123456+dirty", "latest"},
		{"6d5c61e-dirty", "latest"},
		{"6d5c61e", "latest"},
		{"v1.2", "latest"},
		{"1.2.3", "latest"},
		{"dev", "latest"},
		{"", "latest"},
	} {
		if got := installVersion(tc.in); got != tc.want {
			t.Errorf("installVersion(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestWithGitattribute(t *testing.T) {
	const line = ".rhiza/rendered/** linguist-generated=true"
	for _, tc := range []struct {
		name, have, want string
		changed          bool
	}{
		{"Empty", "", line + "\n", true},
		{"Others", "*.go diff=golang\n", "*.go diff=golang\n" + line + "\n", true},
		{"NoFinalNewline", "*.go diff=golang", "*.go diff=golang\n" + line + "\n", true},
		{"KeepsBytes", "a  b\r\n\n# c\r\n", "a  b\r\n\n# c\r\n" + line + "\n", true},
		{"Present", "x\n" + line + "\ny\n", "x\n" + line + "\ny\n", false},
		{"PresentCRLF", line + "\r\n", line + "\r\n", false},
		{"PresentTrailingSpace", line + "  \n", line + "  \n", false},
		{"PresentLastLine", "x\n" + line, "x\n" + line, false},
		{"OtherValue", line + " -diff\n", line + " -diff\n" + line + "\n", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, changed := withGitattribute([]byte(tc.have))
			if string(got) != tc.want || changed != tc.changed {
				t.Errorf("withGitattribute(%q) = %q, %v; want %q, %v", tc.have, got, changed, tc.want, tc.changed)
			}
		})
	}
}

func ciInstall(t *testing.T, ctx *Ctx, req string) CIInstallResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "ci.install", req)
	if err != nil {
		t.Fatal(err)
	}
	return resp.(CIInstallResponse)
}

func readRoot(t *testing.T, ctx *Ctx, p string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(ctx.Root, filepath.FromSlash(p)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func writeRoot(t *testing.T, ctx *Ctx, p, data string) {
	t.Helper()
	p = filepath.Join(ctx.Root, filepath.FromSlash(p))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCIInstall(t *testing.T) {
	both := []string{gitattributesPath, workflowPath}
	ctx := initCtx(t)
	ctx.Version = "v1.2.3"
	r := ciInstall(t, ctx, `{}`)
	if !slices.Equal(r.Written, both) || r.Unchanged == nil || len(r.Unchanged) != 0 {
		t.Fatalf("first install = %+v", r)
	}
	wf := readRoot(t, ctx, workflowPath)
	for _, s := range []string{
		"rhiza ci install",
		"go install github.com/dphaener/rhiza/cmd/rhiza@v1.2.3\n",
		"rhiza check --base origin/${{ github.base_ref }}\n",
		"rhiza check --base ${{ github.event.before }}\n",
		"fetch-depth: 0\n",
		"github.event.before != '" + strings.Repeat("0", 40) + "'\n",
		"github.event.before == '" + strings.Repeat("0", 40) + "'\n      run: rhiza check\n",
	} {
		if !strings.Contains(wf, s) {
			t.Errorf("workflow lacks %q:\n%s", s, wf)
		}
	}
	if got := readRoot(t, ctx, gitattributesPath); got != gitattribute+"\n" {
		t.Errorf(".gitattributes = %q", got)
	}
	if out := human(t, r); out != "Installed the Rhiza check (2 written, 0 unchanged).\n\n- wrote .gitattributes\n- wrote .github/workflows/rhiza.yml\n" {
		t.Errorf("human:\n%s", out)
	}

	r = ciInstall(t, ctx, `{}`)
	if r.Written == nil || len(r.Written) != 0 || !slices.Equal(r.Unchanged, both) {
		t.Errorf("second install = %+v", r)
	}
	if out := human(t, r); out != "Installed the Rhiza check (0 written, 2 unchanged).\n" {
		t.Errorf("human:\n%s", out)
	}

	// A different workflow is a conflict and writes neither file.
	writeRoot(t, ctx, workflowPath, "name: other\n")
	writeRoot(t, ctx, gitattributesPath, "*.go diff=golang\n")
	_, err := callOp(t, ctx, "ci.install", `{}`)
	wantExit(t, err, 3)
	wantNext(t, err, "rhiza ci install --force")
	if got := readRoot(t, ctx, workflowPath); got != "name: other\n" {
		t.Errorf("conflict rewrote the workflow: %q", got)
	}
	if got := readRoot(t, ctx, gitattributesPath); got != "*.go diff=golang\n" {
		t.Errorf("conflict rewrote .gitattributes: %q", got)
	}

	r = ciInstall(t, ctx, `{"force":true}`)
	if !slices.Equal(r.Written, both) || len(r.Unchanged) != 0 {
		t.Errorf("forced install = %+v", r)
	}
	if got := readRoot(t, ctx, workflowPath); got != wf {
		t.Errorf("forced workflow = %q", got)
	}
	if got := readRoot(t, ctx, gitattributesPath); got != "*.go diff=golang\n"+gitattribute+"\n" {
		t.Errorf(".gitattributes = %q", got)
	}
}

func TestCIInstallLatest(t *testing.T) {
	ctx := initCtx(t)
	ctx.Version = "dev"
	ciInstall(t, ctx, `{}`)
	if wf := readRoot(t, ctx, workflowPath); !strings.Contains(wf, "/cmd/rhiza@latest\n") {
		t.Errorf("workflow:\n%s", wf)
	}
}

func TestCIInstallErrors(t *testing.T) {
	ctx := testCtx(t.TempDir())
	_, err := callOp(t, ctx, "ci.install", `{}`)
	wantExit(t, err, 4)
	wantNext(t, err, "rhiza init")

	ctx = initCtx(t)
	ctx.DryRun = true
	_, err = callOp(t, ctx, "ci.install", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza ci install")
	for _, p := range []string{workflowPath, gitattributesPath} {
		if _, serr := os.Stat(filepath.Join(ctx.Root, filepath.FromSlash(p))); serr == nil {
			t.Errorf("a dry run wrote %s", p)
		}
	}
}
