# check-whatif-finish-cheaper

capability · `01M7P2X9DSE7CVKCPE6RRH8AXN` · parent [check](../check.md)

## acceptance

go test -trimpath -count=1 ./internal/model/... ./internal/check/... passes, and make self-check passes. Checks: block gives RHIZA_CPUPROFILE of $TMPDIR/rhiza-x check, built at the parent commit and at the commit, with the pprof -top -cum lines of model.whatIf, model.whatIfLayer and (*builder).finishModel, their share falling; and $TMPDIR/rhiza-base, built at the change set's base 5d2596d (a git worktree under $TMPDIR), and $TMPDIR/rhiza-x, built at the commit, run alternately five times each back to back in the main checkout after one warm call of each, each run's user+sys CPU taken from /usr/bin/time -f '%e s %U user %S sys' for rhiza check, the median of $TMPDIR/rhiza-x's at most 59% (5 s / 8.4 s) of the median of $TMPDIR/rhiza-base's. Failure case: the parent commit's and the commit's binaries' rhiza check --json in the main checkout print the same bytes and exit code, and a new test in internal/model/view_test.go builds a what-if view over a change set whose op the fold refuses and gets the same problem as before the cut; the existing model and check tests keep passing unchanged

## files

internal/model/view.go, internal/model/apply.go (finishModel only), internal/model/view_test.go

## must_not

Change which findings rhiza check reports, their severity, rule or text; change the log's file format; lower the 5 s target or any budget test threshold

## statement

Each in-flight change set's what-if view in rhiza check costs less to finish: model.whatIf's whatIfLayer and (*builder).finishModel (8.4% of the 2026-11-18 profile of rhiza check, beside the Seed.builder copying whatif-copies-only-written cuts) are cut in internal/model, and once check-log-decode-cheaper, whatif-copies-only-written and check-inflight-links-own-only have landed with it, rhiza check's CPU is within m39-check-fast's 5 s target in proportion: at most 5/8.4 of the base binary's

## Links

### Out

- depends-on [whatif-copies-only-written](whatif-copies-only-written.md)
- depends-on [check-inflight-links-own-only](check-inflight-links-own-only.md)
- depends-on [check-log-decode-cheaper](check-log-decode-cheaper.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- verifies [check-whatif-finish-cheaper-verified](../../../../evidence/check-whatif-finish-cheaper-verified.md)
- verifies [check-whatif-finish-cheaper-qa-1](../../../../evidence/check-whatif-finish-cheaper-qa-1.md)
- verifies [check-whatif-finish-cheaper-qa-2](../../../../evidence/check-whatif-finish-cheaper-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
