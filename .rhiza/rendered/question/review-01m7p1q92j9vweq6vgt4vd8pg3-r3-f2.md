# review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f2

question · `01M7P39FWBFB39T4FEQT61HS6F`

## statement

Plan review round 3 left finding f2 (missing-acceptance on inbox-run-list-cached): The acceptance's new test cannot fail before the change. It lists twice over one store and asserts the second call parses only the journal appended between them. internal/run/store.go already has an in-process journalCache, keyed by path, size and mtime, that does exactly this, and TestStoreCachesUnchangedJournals in internal/run/store_cache_test.go already asserts it. The leaf's real change is a cache kept in the user cache folder that outlives the process, which is what makes a second rhiza inbox call cheaper. Nothing in the acceptance checks that: no fresh Store or cleared in-process cache reads the on-disk cache, and no test runs under a test-set user cache dir. The critic confirmed it: internal/run/store.go:263-267 already has a process-wide journalCache keyed by path, size and mtime, used by Store.List (store.go:248-259), and TestStoreCachesUnchangedJournals (internal/run/store_cache_test.go:11) already asserts that an unchanged journal is parsed once and an appended one again. The acceptance's new test (List twice over one store, the second parses only the appended journal) therefore passes before the change. Nothing in the acceptance checks the leaf's real change, a cache in the user cache folder that a fresh process or cleared in-process cache reads. An acceptance that cannot fail before the change: major stands. How should it be settled?

## Links

### Out

- blocks [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)

### In

- answers [m39-inbox-cache-acceptance-on-disk](../decision/m39-inbox-cache-acceptance-on-disk.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
