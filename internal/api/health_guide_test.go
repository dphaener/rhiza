package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestHealthGuide covers the guide part: ok for a current guide or none,
// warn with the install command for an older version, a missing skill, or
// a stale project-scope guide; and its place after usage.
func TestHealthGuide(t *testing.T) {
	home := t.TempDir()
	ctx := schedHealthCtx(t, "")
	ctx.Home = home
	ctx.Version = "v1.2.3"
	scopes := []inboxScope{{"", ctx}}
	homeDir := filepath.Join(home, ".claude", "skills", "rhiza")

	if p := healthGuide(ctx, scopes); p.State != HealthOK || p.Part != HealthGuide {
		t.Fatalf("no guide = %+v", p)
	}

	guideInstall(t, ctx, `{"agent":"claude"}`)
	if p := healthGuide(ctx, scopes); p.State != HealthOK {
		t.Fatalf("current guide = %+v", p)
	}

	ctx.Version = "v1.2.4"
	p := healthGuide(ctx, scopes)
	if p.State != HealthWarn || p.Next != "rhiza guide install --agent claude" || len(p.Detail) != 1 || !strings.Contains(p.Detail[0], homeDir) {
		t.Fatalf("older guide = %+v", p)
	}

	guideInstall(t, ctx, `{"agent":"claude"}`)
	if err := os.Remove(filepath.Join(homeDir, filepath.FromSlash(pluginBuilderPath))); err != nil {
		t.Fatal(err)
	}
	if p := healthGuide(ctx, scopes); p.State != HealthWarn || p.Next != "rhiza guide install --agent claude" {
		t.Fatalf("guide missing skills/rhiza-builder = %+v", p)
	}

	guideInstall(t, ctx, `{"agent":"claude"}`)
	guideInstall(t, ctx, `{"agent":"claude","project":true}`)
	if p := healthGuide(ctx, scopes); p.State != HealthOK {
		t.Fatalf("both current = %+v", p)
	}
	ctx.Version = "v1.2.5"
	guideInstall(t, ctx, `{"agent":"claude"}`)
	p = healthGuide(ctx, scopes)
	if p.State != HealthWarn || p.Next != "rhiza guide install --agent claude --project" {
		t.Fatalf("stale project guide = %+v", p)
	}
}

func TestHealthGuideAfterUsage(t *testing.T) {
	ctx := schedHealthCtx(t, "")
	ctx.Home = t.TempDir()
	resp, err := callOp(t, ctx, "health", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(HealthResponse)
	var order []string
	for _, p := range r.Parts {
		order = append(order, p.Part)
	}
	u, g := -1, -1
	for i, n := range order {
		switch n {
		case HealthUsage:
			u = i
		case HealthGuide:
			g = i
		}
	}
	if u < 0 || g != u+1 {
		t.Fatalf("parts = %v", order)
	}
	if r.ExitCode() != 0 {
		t.Errorf("exit = %d, state %s", r.ExitCode(), r.State)
	}
}
