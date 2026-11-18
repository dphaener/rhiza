# inbox-proposals-overlay-once

capability · `01M7P20XEZF2EG09V98GMRZZ94` · parent [inbox](../inbox.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox|Proposal' passes, with a new test in internal/api/proposals_test.go counting the overlays built for an inbox over several proposed and draft change sets, some proposed sets covered by a run: one overlay per proposed set with no covering run, none for a covered proposed set or a draft set, and none when every proposed set has a run. Checks: block gives RHIZA_CPUPROFILE of $TMPDIR/rhiza-x inbox in the main checkout (warm) before and after, showing inboxUncoveredProposals' share fall, and /usr/bin/time -f '%e s' of the warm call both times. Failure case: a proposed change set whose what-if has a conflict still gets its accept item with the same text as before (the existing proposals tests keep passing unchanged)

## files

internal/api/proposals.go, internal/api/proposals_test.go, internal/api/inbox.go only where it builds the fold the proposals read

## must_not

Change which proposals the inbox lists or their item text; change internal/index or internal/model; change the --json shape of rhiza inbox

## statement

rhiza inbox's uncovered-proposals items (inboxUncoveredProposals in internal/api/proposals.go) cost only the proposed change sets they report: the proposed sets a run already covers are filtered out (uncoveredProposals) before any view is built, so the what-if overlay the profile shows under writeIndex.buildModel and index.WhatIfOverlayOver is built once for each proposed set that lacks a covering run, in that set's own view, and never for a covered, draft or accepted set

## Links

### Out

- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [inbox-blockers-cheaper](inbox-blockers-cheaper.md)
- depends-on [inbox-rest-cut](inbox-rest-cut.md)
- verifies [inbox-proposals-overlay-once-verified](../../../../evidence/inbox-proposals-overlay-once-verified.md)
- verifies [inbox-proposals-overlay-once-qa-1](../../../../evidence/inbox-proposals-overlay-once-qa-1.md)
- verifies [inbox-proposals-overlay-once-qa-2](../../../../evidence/inbox-proposals-overlay-once-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
