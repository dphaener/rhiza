# check-inflight-links-own-only

capability · `01M7P20XPG0C2NH160Z2RHJCGS` · parent [check](../check.md)

## acceptance

go test -trimpath -count=1 ./internal/check/... passes, with a new test in internal/check giving the same in-flight link findings as before for a model with a bad link created by one in-flight set and a bad as-is link (the as-is one reported by the link rule only, once). Checks: block gives RHIZA_CPUPROFILE of $TMPDIR/rhiza-x check in the main checkout before and after, showing checker.inFlightLinks' share fall, and /usr/bin/time -f '%e s' both times; make self-check passes. Failure case: a link created by an in-flight change set whose endpoint is retired or re-kinded by another op of that same set is still reported

## files

internal/check/check.go, internal/check/check_test.go

## must_not

Change which findings rhiza check reports, their severity, rule or text; change internal/model

## statement

rhiza check's in-flight link rule (checker.inFlightLinks in internal/check/check.go) checks only the links each draft, proposed or accepted change set itself created in its what-if view, instead of running badLinks over every link of every view and keeping those afterwards

## Links

### Out

- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [check-whatif-finish-cheaper](check-whatif-finish-cheaper.md)
- verifies [check-inflight-links-own-only-verified](../../../../evidence/check-inflight-links-own-only-verified.md)
- verifies [check-inflight-links-own-only-qa-1](../../../../evidence/check-inflight-links-own-only-qa-1.md)
- verifies [check-inflight-links-own-only-qa-2](../../../../evidence/check-inflight-links-own-only-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
