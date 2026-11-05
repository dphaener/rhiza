# go-build-cache-sweep-verified

evidence · `01M6NSR3W4EW9AKK62D6DKGT91` · actor [claude-code-9c4c6637](../actor/claude-code-9c4c6637.md)

## result

PASS

## statement

At d5ab5678 ran go test -trimpath -count=1 -v ./internal/gosweep/...: all 4 tests pass, plus go vet and gofmt -l clean. TestSweepDirsRemovesOnlyStale builds a fake GOCACHE and TMPDIR with os.Chtimes back-dated (4h, 5h) and fresh (1h) shard entries, a 9h-old top-level README, an old (3h) and a fresh (30m) go-build* folder and an old other-old folder; SweepDirs removes only the two stale shard entries and the old go-build111 and reports Freed{CacheEntries:2, CacheBytes:30, TempFolders:1, TempBytes:100}. TestCacheOffTouchesNoEntry sets GOCACHE=off: CacheDir() is empty and no cache entry is removed, though the old temp folder still goes. Read the diff: the cutoffs are 3h for the cache and 2h for temp folders; the cache dir is $GOCACHE, else os.UserCacheDir()/go-build; only entries in two-hex shard folders are removed, so trim.txt, README and fuzz/ stay; the temp sweep covers only go-build* folders directly under os.TempDir(). This matches the leaf statement and the go-cache-sweep-shape decision.

## Links

### Out

- verifies [go-build-cache-sweep](../capability/rhiza/orchestration/daemon/go-build-cache-sweep.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
