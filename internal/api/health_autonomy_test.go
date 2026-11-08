package api

import (
	"strings"
	"testing"
)

func TestHealthAutonomyRisk(t *testing.T) {
	const phrase = "outside text reaches agents with no permission checks"

	t.Run("flagged plan warns", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		p := healthAutonomy(ctx, exposureSettings("auto", "owner", "owner", "bypassPermissions"), nil)
		if p.State != HealthWarn {
			t.Fatalf("state %q, want warn", p.State)
		}
		if !strings.Contains(p.Summary, phrase) {
			t.Errorf("summary %q lacks %q", p.Summary, phrase)
		}
		// Exposure lines only: each workflow's autonomy declarations are
		// detail lines too.
		var lines []string
		for _, d := range p.Detail {
			if strings.HasPrefix(d, "workflow ") && strings.Contains(d, "runs from a queued request") {
				lines = append(lines, d)
			}
		}
		if len(lines) != 1 || !strings.Contains(lines[0], "workflow plan runs from a queued request") || !strings.Contains(lines[0], "interview") {
			t.Errorf("detail %q, want one line naming plan, the request and its steps", p.Detail)
		}
		if p.Next != "rhiza guide workflows" {
			t.Errorf("next %q", p.Next)
		}
	})

	t.Run("nothing flagged is ok", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		p := healthAutonomy(ctx, exposureSettings("auto", "auto", "auto", "acceptEdits"), nil)
		if p.State != HealthOK || p.Next != "" || strings.Contains(p.Summary, phrase) {
			t.Errorf("got %+v, want ok and unchanged", p)
		}
		for _, d := range p.Detail {
			if strings.Contains(d, "runs from a queued request") {
				t.Errorf("unexpected detail %q", d)
			}
		}
	})

	t.Run("flagged while paused keeps the paused next", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		s := append(exposureSettings("auto", "owner", "owner", "bypassPermissions"), Setting{Key: "rhiza.daemon.paused", Value: "true"})
		p := healthAutonomy(ctx, s, nil)
		if p.State != HealthWarn || p.Next != "rhiza daemon set rhiza.daemon.paused false" {
			t.Errorf("got %+v", p)
		}
		if !strings.HasPrefix(p.Summary, "daemon paused; ") || !strings.Contains(p.Summary, phrase) {
			t.Errorf("summary %q", p.Summary)
		}
	})
}
