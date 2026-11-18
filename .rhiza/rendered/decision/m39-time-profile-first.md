# m39-time-profile-first

decision · `01M7P20WZ6PD3HASF9AGKQ8YFA`

## reasoning

The change set's intent says to profile each before changing anything and never to lower a target. Package wall times in one make test overlap (a package's time includes the others' load) and the shared VM's load swings, so a leaf is judged on its own before/after CPU on the same box, measured inside its own commit

## statement

Each test-time leaf of 01M7P1Q9 starts from a profile and ends with the same one: its commit's Checks: block gives the before and after measurement (go test -cpuprofile with /usr/bin/time for a package, RHIZA_CPUPROFILE with one scratch-built binary for a command) and the pprof -top -cum lines of the cost it cuts. The cut is the top cost that profile names inside the leaf's statement; no milestone target, budget test threshold, rapid check count or test is lowered, skipped or deleted to get there

## Links

### Out

- derived-from [m39-test-time-research](../source/m39-test-time-research.md)

### In

- decided-by [api-tests-copy-prebuilt-index](../capability/rhiza/integrity/self-check/api-tests-copy-prebuilt-index.md)
- decided-by [scripts-top-cost-cut](../capability/rhiza/integrity/self-check/scripts-top-cost-cut.md)
- decided-by [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)
- decided-by [inbox-proposals-overlay-once](../capability/rhiza/orchestration/inbox/inbox-proposals-overlay-once.md)
- decided-by [whatif-copies-only-written](../capability/rhiza/integrity/check/whatif-copies-only-written.md)
- decided-by [check-inflight-links-own-only](../capability/rhiza/integrity/check/check-inflight-links-own-only.md)
- decided-by [check-log-decode-cheaper](../capability/rhiza/integrity/check/check-log-decode-cheaper.md)
- decided-by [check-whatif-finish-cheaper](../capability/rhiza/integrity/check/check-whatif-finish-cheaper.md)
- decided-by [inbox-blockers-cheaper](../capability/rhiza/orchestration/inbox/inbox-blockers-cheaper.md)
- decided-by [cold-run-ship-scripts-cut](../capability/rhiza/integrity/self-check/cold-run-ship-scripts-cut.md)
- decided-by [inbox-rest-cut](../capability/rhiza/orchestration/inbox/inbox-rest-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
