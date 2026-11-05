# go-cache-sweep-shape

decision · `01M6NRBD14M482Q2BXT8XNQ1W5`

## reasoning

I took the asker's recommendations on all three points. Intent: the sibling fix test-temp-folders-swept, a disk-hygiene cleanup like this one, served milestone-25. milestone-12 ('The night runs clean') is about holding work when disk is low, not about freeing space. Timing: Go refreshes a cache entry's mtime only when it is over an hour old, so a 3-hour cutoff never removes an entry a live build is using. Go builds put go-build* work folders in TMPDIR and remove them when they finish, so one older than 2 hours was left by a killed build. That fits the change set's intent ('untouched for a few hours', 'older than a couple of hours'). An hourly sweep plus an early sweep at 2×minFree acts before the existing minFree guard (internal/api/daemon.go daemonDiskLow, default 10G) holds new work, and the 10-minute floor stops repeated scans of a 47G cache. Discovery: resolving $GOCACHE or os.UserCacheDir()/go-build is what go itself does by default, and it avoids exec-ing go from the daemon. os.TempDir() honours TMPDIR, so the daemon must run with the same TMPDIR as its agents (~/tmp on the VM). The setting follows the existing rhiza.daemon.* keys in internal/api/daemon_settings.go, and logging only when something was freed matches the owner's rule that information goes to logs, not to watcher panes.

## statement

Change set 01M6NR6E serves milestone-25 (the backlog). The daemon sweeps once when it starts, then hourly, and also at once when free space drops under 2×rhiza.daemon.minFree, at most once every 10 minutes. Each sweep removes Go build cache entries whose mtime is over 3 hours old and go-build* folders directly under os.TempDir() whose mtime is over 2 hours old. The cache is $GOCACHE, or os.UserCacheDir()/go-build when GOCACHE is unset (no exec of go env); GOCACHE=off skips the cache sweep. A new setting rhiza.daemon.trimGoCache (default on) turns the sweep off. The daemon writes one stderr log line saying what was freed, only when something was freed.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-4c7d9954-1](../question/question-4c7d9954-1.md)
- derived-from [go-cache-sweep-plan](../source/go-cache-sweep-plan.md)

### In

- decided-by [go-build-cache-sweep](../capability/rhiza/orchestration/daemon/go-build-cache-sweep.md)
- decided-by [daemon-trims-go-cache](../capability/rhiza/orchestration/daemon/daemon-trims-go-cache.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
