package api

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// tmpLowCtx is a test ctx under RHIZA_TEST_NOW whose temp free space is
// tmp (unset when empty).
func tmpLowCtx(t *testing.T, tmp string) (*Ctx, *bytes.Buffer) {
	t.Helper()
	ctx := testCtx(t, t.TempDir())
	var log bytes.Buffer
	ctx.Stderr = &log
	prev := ctx.Getenv
	ctx.Getenv = func(k string) string {
		switch k {
		case "RHIZA_TEST_NOW":
			return "2026-10-24T12:00:00Z"
		case "RHIZA_TEST_TMP_FREE_BYTES":
			return tmp
		}
		return prev(k)
	}
	return ctx, &log
}

// Under RHIZA_TEST_NOW the temp free space comes only from
// RHIZA_TEST_TMP_FREE_BYTES: without it the figure is not known and nothing
// holds, and the host's temp space is never read (diskFreeBytes panics in
// this package's tests).
func TestTmpFreeBytesPinnedUnderTest(t *testing.T) {
	t.Parallel()
	ctx, _ := tmpLowCtx(t, "")
	if _, err := tmpFreeBytes(ctx); err == nil {
		t.Fatal("tmpFreeBytes under RHIZA_TEST_NOW without RHIZA_TEST_TMP_FREE_BYTES read a figure")
	}
	if low := daemonTmpLow(ctx, 2<<30, nil); low != nil {
		t.Fatalf("unknown temp space held: %+v", low)
	}
	ctx, _ = tmpLowCtx(t, "1073741824")
	if free, err := tmpFreeBytes(ctx); err != nil || free != 1<<30 {
		t.Fatalf("tmpFreeBytes = %d, %v; want 1 GiB", free, err)
	}
}

// Free temp space under the minimum holds with the free space and the
// minimum, keeps its start while it lasts, and ends at the minimum; a
// minimum of 0 never holds. Each hold and release is logged once.
func TestTmpLowHold(t *testing.T) {
	t.Parallel()
	ctx, log := tmpLowCtx(t, "1073741824")
	low := daemonTmpLow(ctx, 2<<30, nil)
	if low == nil || low.Free != 1<<30 || low.Min != 2<<30 || !low.Since.Equal(ctx.Now().UTC()) {
		t.Fatalf("hold = %+v, want 1 GiB free under 2 GiB since now", low)
	}
	if got, want := tmpLowDetail(low), "1.0 GiB free in the temp directory, under the minimum of 2.0 GiB"; got != want {
		t.Fatalf("detail %q, want %q", got, want)
	}
	if !strings.Contains(log.String(), "holding new runs: tmp-low: 1.0 GiB free") {
		t.Fatalf("log %q: want the hold", log.String())
	}
	earlier := &daemon.DiskLow{Since: ctx.Now().Add(-time.Hour), Free: 1 << 30, Min: 2 << 30}
	log.Reset()
	if again := daemonTmpLow(ctx, 2<<30, earlier); again == nil || !again.Since.Equal(earlier.Since) {
		t.Fatalf("lasting hold = %+v, want since %v", again, earlier.Since)
	}
	if log.Len() != 0 {
		t.Fatalf("lasting hold logged %q", log.String())
	}
	if off := daemonTmpLow(ctx, 0, earlier); off != nil {
		t.Fatalf("minTmpFree 0 held: %+v", off)
	}
	if !strings.Contains(log.String(), "released the tmp-low hold on new runs") {
		t.Fatalf("log %q: want the release", log.String())
	}
	ctx, _ = tmpLowCtx(t, "2147483648")
	if at := daemonTmpLow(ctx, 2<<30, nil); at != nil {
		t.Fatalf("free at the minimum held: %+v", at)
	}
}

// The hold's inbox item and its skip both name the free space and the
// minimum; the load and memory hold goes first.
func TestTmpLowItemAndSkip(t *testing.T) {
	t.Parallel()
	low := &daemon.DiskLow{Since: time.Date(2030, 1, 1, 0, 1, 0, 0, time.UTC), Free: 1 << 30, Min: 2 << 30}
	items := tmpLowItems(low)
	if len(items) != 1 || items[0].Kind != InboxTmpLow || items[0].At != "2030-01-01T00:01:00Z" ||
		items[0].Detail != "1.0 GiB free in the temp directory, under the minimum of 2.0 GiB: the daemon holds new starts and resumes" {
		t.Fatalf("items %+v", items)
	}
	if !inboxMachineWide(InboxTmpLow) {
		t.Fatal("tmp-low is not machine-wide")
	}
	if tmpLowItems(nil) != nil || tmpLowHold(nil) != nil {
		t.Fatal("no hold made an item or a hold")
	}
	h := daemonConf{tmpLow: low}.hostHold()
	if h == nil || h.Reason != daemon.SkipTmpLow || h.Detail != tmpLowDetail(low) {
		t.Fatalf("host hold %+v", h)
	}
	load := &daemon.MachineHold{Reason: daemon.SkipLoadHigh}
	if got := (daemonConf{machine: load, tmpLow: low}).hostHold(); got != load {
		t.Fatalf("host hold %+v, want the load hold", got)
	}
}

// A tmp-low hold is a pre-start wait: the queue-wait journal records it
// with its detail, and the change set starts once it ends.
func TestQueueWaitTmpLow(t *testing.T) {
	t.Parallel()
	f := usageFix(t)
	hold := tmpLowHold(&daemon.DiskLow{Since: f.clock, Free: 1 << 30, Min: 2 << 30})
	if started := queueWaitPoll(t, f, 4, nil, hold, false); len(started) != 0 {
		t.Fatalf("started %v under a tmp-low hold", started)
	}
	got := wantQueueWaits(t, f, "tmp-low")
	if got[0].Detail != "1.0 GiB free in the temp directory, under the minimum of 2.0 GiB" {
		t.Fatalf("detail %q", got[0].Detail)
	}
	if started := queueWaitPoll(t, f, 4, nil, nil, false); len(started) != 1 || started[0] != f.target {
		t.Fatalf("started %v", started)
	}
	wantQueueWaits(t, f, "tmp-low", "started")
}
