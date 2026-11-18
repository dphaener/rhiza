# inbox-proposals-overlay-once-verified

evidence · `01M7PMF8PADWBQ4FYXTJ50XYVZ` · actor [claude-code-2a4a1ce0](../actor/claude-code-2a4a1ce0.md)

## result

PASS

## statement

Verified at 493aa35f7: go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox|Proposal' exit 0 (12.3s), new TestInboxUncoveredProposalsOverlayOnce counts views via inboxFold.onProposalView over 4 proposed + 2 draft sets (no runs: 4 views; 2 covered: 2 views; all covered: 0), existing proposals tests unchanged and passing. Before/after binaries (HEAD~1 vs HEAD) in the main checkout, warm, back to back at load ~4: RHIZA_CPUPROFILE inboxUncoveredProposals 31.33%->16.28% and 33.82%->9.80%, WhatIfOverlayOver 28.92%->11.63% and 32.35%->3.92%; /usr/bin/time profiled 0.73->0.52 s, 0.62->0.53 s; unprofiled 0.54->0.46 s, 0.69->0.32 s. rhiza inbox --json before/after byte-identical (cmp). Diff touches only proposals.go, proposals_test.go and inbox.go's fold struct; no internal/index or internal/model change. make self-check in a clean detached worktree at HEAD exit 0.

## Links

### Out

- verifies [inbox-proposals-overlay-once](../capability/rhiza/orchestration/inbox/inbox-proposals-overlay-once.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
