# daemon-trims-go-cache

capability · `01M6NRETPP5SJ8RMGV901X4G33` · parent [daemon](../daemon.md)

## acceptance

A unit test of the sweep schedule with an injected clock and free space shows a sweep is due at start, again after an hour, at once under 2x minFree and not again within 10 minutes; a daemon test with a fake GOCACHE holding a back-dated entry shows daemon --once removes it and logs the freed line, and with rhiza.daemon.trimGoCache false leaves it; daemon settings lists rhiza.daemon.trimGoCache with default true; describe-golden.txtar is regenerated

## statement

The daemon runs the Go build cache sweep when it starts, then hourly, and at once when free space drops under 2x rhiza.daemon.minFree (at most once every 10 minutes), unless rhiza.daemon.trimGoCache is off, writing one stderr line with what was freed only when something was

## Links

### Out

- decided-by [go-cache-sweep-shape](../../../../decision/go-cache-sweep-shape.md)

### In

- verifies [daemon-trims-go-cache-verified](../../../../evidence/daemon-trims-go-cache-verified.md)
- verifies [daemon-trims-go-cache-qa-1](../../../../evidence/daemon-trims-go-cache-qa-1.md)
- verifies [pr-303-the-daemon-keeps-the-go-build](../../../../evidence/pr-303-the-daemon-keeps-the-go-build.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
