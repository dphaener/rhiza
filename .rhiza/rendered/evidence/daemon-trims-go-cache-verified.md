# daemon-trims-go-cache-verified

evidence · `01M6NTD248JZ0ADGFX090PW12C` · actor [claude-code-0073b79b](../actor/claude-code-0073b79b.md)

## result

PASS

## statement

At 230826ed: go test -trimpath ./internal/daemon/... ./internal/gosweep/... pass; TestSweepSchedule (injected clock and free space) shows a sweep due at start, not at 5/59 min, due after an hour, urgent at 10 min under 2x minFree but not at 9 min, not at exactly 2x, not with guard off or unknown free space, and not within 10 min of an urgent sweep. TestScripts/daemon-trims-go-cache passes: with a fake GOCACHE holding an entry back-dated 6h, daemon --once with rhiza.daemon.trimGoCache false keeps it and logs nothing; set true it removes it, keeps the fresh entry and logs one 'go build sweep freed ...: 1 cache entries, 0 temp folders' line; a third run logs nothing; daemon settings lists 'rhiza.daemon.trimGoCache = true'. describe-golden passes and -update leaves it unchanged (settings keys are not in the schema), so it is current. All TestScripts/daemon-* pass. Read the diff: sweep runs per poll gated by conf.trimGoCache, GOCACHE/TMPDIR read via ctx env. Nit only: chtimes was inserted between foldOK's doc comment and foldOK in cmd/rhiza/main_test.go, so foldOK's comment now sits above chtimes.

## Links

### Out

- verifies [daemon-trims-go-cache](../capability/rhiza/orchestration/daemon/daemon-trims-go-cache.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
