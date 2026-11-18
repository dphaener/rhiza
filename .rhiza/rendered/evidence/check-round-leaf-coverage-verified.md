# check-round-leaf-coverage-verified

evidence · `01M7PY1560AF6WEENB1XNA6T6N` · actor [claude-code-c5405785](../actor/claude-code-c5405785.md)

## result

PASS

## statement

At b1483bb53: go test -trimpath -p 2 ./internal/api/... ok (177s) and ./internal/run/... ok; TestRunCheckFindingsLeafCoverage covers each acceptance point (evidence per leaf journaled on entry and round; finding on leaf-a + evidence on leaf-b accepted; no-leaf change set and classic review step accept {"findings":[]}; conductor {"findings":[]} refused naming leaf-a and leaf-b, Next=runFindingsNext, journal unchanged; non-leaf, unknown, blank-line and no-leaf evidence refused). go test ./cmd/rhiza -run 'TestScripts/describe-golden$' ok; make self-check exit 0 in a clean detached worktree at HEAD (it fails in the run worktree only because of untracked session actor records). ship-change-one-ci exits 1 as the Checks block predicts: its fake conductor's empty round is now refused, and the script is unchanged (it belongs to leaf conduct-fixtures-cover-leaves). Diff: the rule is keyed on step.Conducts(), not check names; it has no time heuristic and uses live (non-retired) leaves; review-check, --leaf and plan-review decoding are unchanged (DisallowUnknownFields still refuses evidence there); describe-golden.txtar is the only testscript edited; no TODOs.

## Links

### Out

- verifies [check-round-leaf-coverage](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-leaf-coverage.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
