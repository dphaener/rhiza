package api

import (
	"strings"
	"testing"
)

// pinnedFreeBytes is the free disk figure every API test sees unless it sets
// its own: 1 TiB, far over the daemon's 10G default minimum.
const pinnedFreeBytes = "1099511627776"

// pinFreeBytes makes ctx see pinnedFreeBytes of free disk, beside
// RHIZA_TEST_NOW, keeping whatever Getenv it already has. Tests reach
// repoFreeBytes only through a ctx so pinned: the real read fails the test.
func pinFreeBytes(ctx *Ctx) {
	prev := ctx.Getenv
	ctx.Getenv = func(k string) string {
		switch k {
		case "RHIZA_TEST_NOW":
			return "2026-11-08T10:30:00Z"
		case "RHIZA_TEST_FREE_BYTES":
			return pinnedFreeBytes
		}
		if prev != nil {
			return prev(k)
		}
		return ""
	}
}

// setFreeBytes overrides the pinned figure on ctx with its own.
func setFreeBytes(ctx *Ctx, n string) {
	prev := ctx.Getenv
	ctx.Getenv = func(k string) string {
		if k == "RHIZA_TEST_FREE_BYTES" {
			return n
		}
		return prev(k)
	}
}

// A test can still set its own free figure: under 10 GiB the daemon's disk
// guard reports low and health's disk part goes bad.
func TestFreeBytesOverride(t *testing.T) {
	f := newSessionFixture(t, "s")
	pinFreeBytes(f.ctx)
	if low, free := daemonDiskLow(f.ctx, 10<<30); low || free != 1<<40 {
		t.Fatalf("pinned: low=%v free=%d, want not low at 1 TiB", low, free)
	}
	setFreeBytes(f.ctx, "5368709120")
	if low, free := daemonDiskLow(f.ctx, 10<<30); !low || free != 5<<30 {
		t.Fatalf("override: low=%v free=%d, want low at 5 GiB", low, free)
	}
	part := healthDisk(f.ctx, []Setting{{Key: "rhiza.daemon.minFree", Value: "10G"}}, nil)
	if part.State != HealthBad || !strings.Contains(part.Summary, "under the minimum") {
		t.Fatalf("health disk part = %+v, want bad", part)
	}
}
