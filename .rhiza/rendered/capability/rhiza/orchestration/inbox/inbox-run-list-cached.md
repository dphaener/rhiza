# inbox-run-list-cached

capability · `01M7P20XBYDN7HR4ZRG7W101HX` · parent [inbox](../inbox.md)

## acceptance

go test -trimpath -count=1 ./internal/run passes, with a new test in internal/run/store_test.go that can only pass with the on-disk cache: with t.Setenv("XDG_CACHE_HOME", t.TempDir()) and t.Setenv("HOME", t.TempDir()) so no test reads the host's user cache, List once over a store, append to one journal, then clear the process-wide journalCache and List again over a fresh Store; journalReads counts only the appended journal, the result equals an uncached List, and the cache file lies under the test-set cache dir, not the repository or .rhiza. Checks: block gives RHIZA_CPUPROFILE=$TMPDIR/inbox.prof $TMPDIR/rhiza-x inbox in the main checkout (one warm call first) before and after, with go tool pprof -top -cum showing run.(*Store).List's share fall, and /usr/bin/time -f '%e s' of the warm call both times. go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox' passes. Failure case, same setup (test-set XDG_CACHE_HOME and HOME, process-wide journalCache cleared, fresh Store): a cache file that is corrupt or from another format version is ignored and every journal is parsed (journalReads equals the journal count), with List equal to an uncached List (tested in internal/run/store_test.go)

## files

internal/run/store.go, internal/run/store_test.go, internal/api/inbox.go only where it calls rs.List

## must_not

Change the run journal's file format or what List returns; write the cache under the repository or .rhiza; read the host's real user cache in a test; serve a changed journal from the cache

## statement

rhiza inbox reads the run journals that changed since its last call and reuses, from a cache in the user cache folder keyed by each journal's path, size and modification time, what it parsed of the rest, so run.Store.List's share of an inbox call falls to the journals that changed

## Links

### Out

- decided-by [m39-inbox-run-summaries](../../../../decision/m39-inbox-run-summaries.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)
- decided-by [m39-inbox-cache-acceptance-on-disk](../../../../decision/m39-inbox-cache-acceptance-on-disk.md)

### In

- depends-on [inbox-blockers-cheaper](inbox-blockers-cheaper.md)
- depends-on [inbox-rest-cut](inbox-rest-cut.md)
- blocks [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f2](../../../../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f2.md)
- verifies [inbox-run-list-cached-verified](../../../../evidence/inbox-run-list-cached-verified.md)
- verifies [inbox-run-list-cached-qa-1](../../../../evidence/inbox-run-list-cached-qa-1.md)
- verifies [inbox-run-list-cached-qa-2](../../../../evidence/inbox-run-list-cached-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
