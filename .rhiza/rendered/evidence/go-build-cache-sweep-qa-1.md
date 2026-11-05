# go-build-cache-sweep-qa-1

evidence · `01M6NTQQP57W6PATBPR2XZVWDX` · actor [claude-code-31ee992e](../actor/claude-code-31ee992e.md)

## result

PASS

## statement

QA at e129e9f1 with the built binary (go build -trimpath) in a scratch repo, real clock: a fake GOCACHE held 5h-old (touch -d) entries in shards ab/ and 0f/, a fresh ab/ entry, a 5h-old entry in non-shard zz/ and a 5h-old top-level README; a fake TMPDIR held a 3h-old go-build123 (2000-byte file), a 30m-old go-build999 and a 3h-old other-old. rhiza daemon --once removed only the two stale shard entries and go-build123, logging 'go build sweep freed 6.8 KiB: 2 cache entries, 1 temp folders' (5000+2+2000 bytes, correct); fresh-a, zz/, README, go-build999 and other-old stayed. With GOCACHE=off a 5h-old ab/ entry stayed and no sweep line was logged. With GOCACHE unset and XDG_CACHE_HOME pointed at a fake dir, the 5h-old entry under $XDG_CACHE_HOME/go-build/cd/ was removed (os.UserCacheDir()/go-build fallback). Scoped go test of internal/gosweep passes.

## Links

### Out

- verifies [go-build-cache-sweep](../capability/rhiza/orchestration/daemon/go-build-cache-sweep.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
