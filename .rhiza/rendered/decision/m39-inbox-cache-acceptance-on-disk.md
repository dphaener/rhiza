# m39-inbox-cache-acceptance-on-disk

decision · `01M7P3C9XKMZ1Q12RWKKD1F8Z5`

## reasoning

The finding holds: internal/run/store.go:248-267 already keeps a process-wide journalCache keyed by path, size and mtime, and TestStoreCachesUnchangedJournals (internal/run/store_cache_test.go:11) already asserts an unchanged journal is parsed once, so the acceptance's test passes before the change. The leaf's statement is a cache in the user cache folder that outlives the process, which is what makes a second rhiza inbox call cheaper; a test only proves that if the in-process cache is out of the way (cleared, fresh Store) so the parse saving can only come from the disk. Its must_not already forbids reading the host's real user cache in a test and writing under the repository or .rhiza, so the test sets XDG_CACHE_HOME (as internal/api tests like inbox_gate_budget_test.go already do) and checks where the file lands. Tightening the acceptance keeps the leaf's scope and intent; dropping or rescoping the leaf would lose the milestone's inbox 200 ms target work.

## statement

Fix inbox-run-list-cached's acceptance so its new test can only pass with the on-disk cache: in internal/run/store_test.go, with t.Setenv("XDG_CACHE_HOME", t.TempDir()) (and HOME) so no test reads the host's user cache, List once over a store, append to one journal, then clear the process-wide journalCache and List again over a fresh Store; journalReads must count only the appended journal, the result must equal an uncached List, and the cache file must lie under the test-set cache dir, not the repository or .rhiza. The failure case keeps the same setup: a corrupt cache file, or one from another format version, is ignored and every journal is parsed (journalReads equals the journal count), with List equal to an uncached List. The rest of the acceptance (the pprof and /usr/bin/time Checks: block, the internal/api Inbox tests) stays as written.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f2](../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f2.md)

### In

- decided-by [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
