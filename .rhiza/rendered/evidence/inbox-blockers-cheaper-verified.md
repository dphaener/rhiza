# inbox-blockers-cheaper-verified

evidence · `01M7PN30KG9CD21MNA9HGA5NC5` · actor [claude-code-e28d3814](../actor/claude-code-e28d3814.md)

## result

PASS

## statement

At d97acc451: go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox|Plan' ok 13.8s (incl. new TestInboxBlockersListProposedNotDraft: proposed set's store-q/own-q listed with their text, draft's cache-q unlisted); -run 'Ready|Roadmap|Question' ok; cmd/rhiza TestScripts/(inbox-blocking-questions|inbox-blocker-advice) ok with RHIZA_/SHIP_/REX_ unset; go vet and gofmt clean. Warm RHIZA_CPUPROFILE of $TMPDIR rhiza inbox --json in the main checkout, parent 493aa35f7 then commit, back to back, 6 runs each merged, load ~4.6: pair 1 inboxBlockers 14.29%->5.24%, planGather 9.41%->planGatherWith 4.80%, planCore 0 samples both; pair 2 inboxBlockers 16.54%->6.17%, planGather 10.90%->5.76%, planCore 0.75%->0 (planBlockersOnly 0 samples); Analysis.graph/Closure 0.7-1.1% -> none. Parent and commit binaries' rhiza inbox --json byte-identical (cmp) in the main checkout and in the run worktree. Diff touches only the leaf's files plus its actor record; no stubs.

## Links

### Out

- verifies [inbox-blockers-cheaper](../capability/rhiza/orchestration/inbox/inbox-blockers-cheaper.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
