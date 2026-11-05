# pr-303-the-daemon-keeps-the-go-build

evidence · `01M6NW7WEXRKQJMNM619BGDP92` · actor [dphaener](../actor/dphaener.md)

## head

da9537bc45e79ca03c5181ea51e503bd15141253

## pr

https://github.com/dphaener/rhiza/pull/303

## result

1 check passed: check SUCCESS

## statement

Pull request #303 (https://github.com/dphaener/rhiza/pull/303) passed CI at da9537b

## status

verified

## Links

### Out

- verifies [go-build-cache-sweep](../capability/rhiza/orchestration/daemon/go-build-cache-sweep.md)
- verifies [daemon-trims-go-cache](../capability/rhiza/orchestration/daemon/daemon-trims-go-cache.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
