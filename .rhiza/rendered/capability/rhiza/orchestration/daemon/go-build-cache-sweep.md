# go-build-cache-sweep

capability · `01M6NREHKHFKZWMXYWG2F38AM5` · parent [daemon](../daemon.md)

## acceptance

A Go test with a fake GOCACHE and TMPDIR holding back-dated (os.Chtimes) and fresh cache entries, an old and a fresh go-build* folder and an old folder with another prefix shows only the stale cache entries and the old go-build* folder go, with the bytes and counts freed reported; with GOCACHE=off no cache entry is touched

## statement

A sweep removes Go build cache entries whose mtime is over 3 hours old, from $GOCACHE or os.UserCacheDir()/go-build when GOCACHE is unset (none when GOCACHE=off), and go-build* folders directly under os.TempDir() whose mtime is over 2 hours old, and reports the bytes and entries it freed

## Links

### Out

- decided-by [go-cache-sweep-shape](../../../../decision/go-cache-sweep-shape.md)

### In

- verifies [go-build-cache-sweep-verified](../../../../evidence/go-build-cache-sweep-verified.md)
- verifies [go-build-cache-sweep-qa-1](../../../../evidence/go-build-cache-sweep-qa-1.md)
- verifies [pr-303-the-daemon-keeps-the-go-build](../../../../evidence/pr-303-the-daemon-keeps-the-go-build.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
