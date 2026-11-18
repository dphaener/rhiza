# check-whatif-finish-cheaper-verified

evidence · `01M7PKFCZ1YYXCSDJ90TGXN13T` · actor [claude-code-32643571](../actor/claude-code-32643571.md)

## result

PASS

## statement

Verified at b7a419cbe + a7e6ab696 in a clean detached worktree: go test -trimpath -count=1 -p 2 ./internal/model/... ./internal/check/... exit 0; TestWhatIfRefusedOpSameProblem passes at the commit and with the new view_test.go on the parent e96d08bb7 (same problems); make self-check exit 0; rhiza render 0 written. Main checkout: parent vs commit rhiza check --json same bytes (2614), both exit 0. RHIZA_CPUPROFILE check parent->commit: model.whatIf 13.51%->8.96%, whatIfLayer 8.71%->6.12%, (*builder).finishModel 5.26%->3.43%. Five alternating runs after a warm call each, user+sys: base 5d2596d 15.98, 15.49, 14.89, 14.58, 15.16 (median 15.16 s); commit 8.75, 6.40, 6.25, 6.42, 7.70 (median 6.42 s) = 42.3% of base, target <= 59%, load ~5. Diff touches only view.go, apply.go finishModel and view_test.go; Closure output is sorted so the BinarySearch is sound

## Links

### Out

- verifies [check-whatif-finish-cheaper](../capability/rhiza/integrity/check/check-whatif-finish-cheaper.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
