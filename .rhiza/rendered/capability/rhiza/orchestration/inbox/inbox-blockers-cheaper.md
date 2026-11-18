# inbox-blockers-cheaper

capability · `01M7P2X9G4SDMT389QMT2XMQE0` · parent [inbox](../inbox.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox|Plan' passes. Checks: block gives RHIZA_CPUPROFILE of a warm $TMPDIR/rhiza-x inbox, built at the parent commit and at the commit, with the pprof -top -cum lines of inboxBlockers, planGather and planCore, their share falling. Failure case: the parent commit's and the commit's binaries' rhiza inbox --json, run back to back in the main checkout, print the same bytes, and a new test in internal/api lists a proposed change set's blocking question as before with the same text while a draft set's blocker stays unlisted; the existing inbox and plan tests keep passing unchanged

## files

internal/api/inbox.go, internal/api/plan.go, internal/api/inbox_test.go, internal/api/plan_test.go

## must_not

Change which items rhiza inbox lists, their text or the --json shape; change what rhiza plan reports; change the run journal or log file format; lower the 200 ms target or any budget test threshold

## statement

rhiza inbox's blocker items cost less: inboxBlockers' planGather and planCore over the to-be view (changeset.(*Analysis).graph and Closure: 6.7% of a warm RHIZA_CPUPROFILE of rhiza inbox in the main checkout on 2026-11-18) are cut in internal/api, so their share of a warm rhiza inbox falls

## Links

### Out

- depends-on [inbox-run-list-cached](inbox-run-list-cached.md)
- depends-on [inbox-proposals-overlay-once](inbox-proposals-overlay-once.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)
- decided-by [m39-inbox-rest-leaf](../../../../decision/m39-inbox-rest-leaf.md)

### In

- depends-on [inbox-rest-cut](inbox-rest-cut.md)
- blocks [review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f1](../../../../question/review-01m7p1q92j9vweq6vgt4vd8pg3-r3-f1.md)
- verifies [inbox-blockers-cheaper-verified](../../../../evidence/inbox-blockers-cheaper-verified.md)
- verifies [inbox-blockers-cheaper-qa-1](../../../../evidence/inbox-blockers-cheaper-qa-1.md)
- verifies [inbox-blockers-cheaper-qa-2](../../../../evidence/inbox-blockers-cheaper-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
