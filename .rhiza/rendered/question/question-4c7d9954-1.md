# question-4c7d9954-1

question · `01M6NRA08PHMC1XWSWFS5MEZX9`

## session

4c7d9954-929b-4ec4-9dc0-01a1ab8367a9

## statement

Which intent should change set 01M6NR6E (the daemon trims the Go build cache and stale go-build* temp folders) serve? I recommend milestone-25 (the backlog), because the matching test-temp-folders-swept fix served it. milestone-12 ("The night runs clean", which holds the disk guard) is the alternative. / When should the daemon sweep? I recommend once when it starts, then hourly, and also at once when free space drops under 2×rhiza.daemon.minFree (no more than once every 10 minutes). Each sweep removes Go cache entries untouched for over 3 hours (Go refreshes an entry's mtime only when it is over an hour old) and go-build* folders in TMPDIR older than 2 hours. / How should the daemon find the cache, and can the sweep be turned off? I recommend: the cache is $GOCACHE, or os.UserCacheDir()/go-build when that is unset (no exec of go env), and GOCACHE=off skips it. Temp folders are go-build* directly under os.TempDir(). A new setting, rhiza.daemon.trimGoCache (default on), turns the sweep off. One stderr log line reports what was freed, only when something was.

## Links

### In

- answers [go-cache-sweep-shape](../decision/go-cache-sweep-shape.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
