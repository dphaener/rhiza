# inbox-rest-cut

capability · `01M7P3FX70P6S111HSJMKA3B88` · parent [inbox](../inbox.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox' passes, plus go test -trimpath -count=1 for each of internal/index and internal/model the cut touches. Checks: block gives RHIZA_CPUPROFILE of a warm $TMPDIR/rhiza-x inbox built at the parent commit (the three inbox leaves landed) and at the commit, with the pprof -top -cum lines of each cost cut, their share falling; and $TMPDIR/rhiza-base, built at the change set's base 5d2596d (a git worktree under $TMPDIR), and $TMPDIR/rhiza-x, built at the commit, run alternately five times each back to back in the main checkout after one warm call of each, each run's user+sys CPU taken from /usr/bin/time -f '%e s %U user %S sys' for rhiza inbox, the median of $TMPDIR/rhiza-x's at most 31% (200 ms / 645 ms) of the median of $TMPDIR/rhiza-base's. Failure case: the parent commit's and the commit's binaries' rhiza inbox --json, run back to back in the main checkout, print the same bytes, and after a new log entry (rhiza note add in a scratch copy of the checkout) both binaries list the same items, so the cut serves no stale index or view; the existing inbox tests keep passing unchanged

## files

internal/api/inbox.go, internal/api/inbox_test.go, and internal/index or internal/model files only where the fresh profile puts a top cost, with their _test.go

## must_not

Change which items rhiza inbox lists, their text or the --json shape; change the log, index or run journal file format; lower the 200 ms target or any budget test threshold

## statement

A warm rhiza inbox's remaining cost is cut once inbox-run-list-cached, inbox-proposals-overlay-once and inbox-blockers-cheaper have landed: the ~32% those three leave (index open and catch-up, the to-be reader, the other inbox lookups), taken from a fresh RHIZA_CPUPROFILE of rhiza inbox at that point, its top costs cut in internal/api, internal/index or internal/model, so a warm rhiza inbox's CPU is within m39-reads-fast's 200 ms target in proportion: at most 200/645 of the base binary's

## Links

### Out

- depends-on [inbox-run-list-cached](inbox-run-list-cached.md)
- depends-on [inbox-proposals-overlay-once](inbox-proposals-overlay-once.md)
- depends-on [inbox-blockers-cheaper](inbox-blockers-cheaper.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)
- decided-by [m39-inbox-rest-leaf](../../../../decision/m39-inbox-rest-leaf.md)

### In

- verifies [inbox-rest-cut-verified](../../../../evidence/inbox-rest-cut-verified.md)
- verifies [inbox-rest-cut-qa-1](../../../../evidence/inbox-rest-cut-qa-1.md)
- verifies [inbox-rest-cut-qa-2](../../../../evidence/inbox-rest-cut-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
