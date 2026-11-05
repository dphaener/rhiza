# daemon-trims-go-cache-qa-1

evidence · `01M6NTQZRK0TCGZW4HZA6ZFW3Z` · actor [claude-code-31ee992e](../actor/claude-code-31ee992e.md)

## result

PASS

## statement

QA at e129e9f1 with the built binary in a scratch repo: rhiza daemon settings lists 'rhiza.daemon.trimGoCache = true (default)'; daemon set ... maybe is refused (exit 2, must be true or false). With trimGoCache false, daemon --permission-mode manual --once left every stale cache entry and go-build* folder and logged no sweep line; set true, the same run swept them and logged one line; a third --once with nothing stale logged no line. A foreground daemon run for 25s with RHIZA_TEST_FREE_BYTES=15G (under 2x the 10G minFree) swept at start (4 entries), and two entries back-dated 5h after start were still there at stop: no repeat sweep within 10 minutes despite low space. Hourly and urgent-after-10-minute timing not observable in real time; covered by TestSweepSchedule, which passes with internal/daemon, plus TestScripts/daemon-trims-go-cache and describe-golden.

## Links

### Out

- verifies [daemon-trims-go-cache](../capability/rhiza/orchestration/daemon/daemon-trims-go-cache.md)

## History

- `01M6NR6EK065HB06P15DYMB0XX` The daemon keeps the Go build cache and stale build temp folders from filling the disk: on 2026-11-05 the VM's ~/.cache/go-build reached 47G and ~/tmp held stale go-build* folders from killed builds (disk 19G free of 97G) after one evening of runs; Go only trims entries unused for 5 days. When free space nears rhiza.daemon.minFree (or on a schedule), the daemon trims go-build entries untouched for a few hours (mtime, which go refreshes on use) and removes go-build* temp folders older than a couple of hours, logging what it freed; a test with a fake cache dir shows only stale entries go.
