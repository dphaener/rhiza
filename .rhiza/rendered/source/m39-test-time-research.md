# m39-test-time-research

source · `01M7P1XZWPBEW388FYBFYYVZF6`

## quotes

```json
[]
```

## raw_sha256

b49edda65aa41dfbe38cd04abc9b59077943a884c964b73be2c7854bba341ff4

## session

f4b48d65-d95a-4af9-a617-01a1ec1bfb86

## summary

Autonomous plan run (RHIZA_AUTONOMOUS=1, no owner interview) for change set 01M7P1Q9, which carries milestone 39's missed targets: cold make test 400 s wall (target 240 s), rhiza inbox 645 ms (target 200 ms), rhiza check 8.4 s (target 5 s). Research pass, 2026-11-18 on the shared VM (load about 3.6, so absolute times run high): the cold log (perf-data/m39/cold.log at 5d2596d) has internal/api 318 s, cmd/rhiza 254 s, internal/index 113 s, internal/slice 33 s, internal/web 30 s; package times in one make test overlap, so a package's wall time includes contention from the others. internal/api's testCtx (init_test.go) gives every test a fresh index under t.TempDir() since api-tests-on-index (#599), and snapshot copies (fixture_test.go copyTo) carry no index, so every test catches an index up from nothing. A RHIZA_CPUPROFILE of rhiza inbox in the main checkout (1.05 s under load) puts 37% in run.(*Store).List, which reads and parses every run journal on each call (journalCache lives only in the process), and 24% in inboxUncoveredProposals through writeIndex.buildModel and index.WhatIfOverlayOver. A profile of rhiza check (12.3 s under load) puts 30% in checker.model, mostly model.WhatIfFrom for each of the ~64 in-flight change sets, of which Seed.builder's copyNode/copyFields of every node is 20%; 15% in checker.inFlightLinks, which runs badLinks over every link of every what-if view and keeps only the links the set itself created; 11% in store.readEntries. Plan: one leaf per cost, each starting from a profile in its Checks: block, no target lowered and no test, check or budget loosened.

## title

Milestone 39 test-time plan: research pass (autonomous)

## Links

### In

- derived-from [m39-no-index-package-leaf](../decision/m39-no-index-package-leaf.md)
- derived-from [m39-time-profile-first](../decision/m39-time-profile-first.md)
- derived-from [m39-api-index-per-test-copy](../decision/m39-api-index-per-test-copy.md)
- derived-from [m39-inbox-run-summaries](../decision/m39-inbox-run-summaries.md)
- derived-from [m39-whatif-copy-on-write](../decision/m39-whatif-copy-on-write.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
