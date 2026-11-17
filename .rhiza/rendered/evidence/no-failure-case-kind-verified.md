# no-failure-case-kind-verified

evidence · `01M7KPKXYE8WFYY60E2M67XH7F` · actor [claude-code-6408249b](../actor/claude-code-6408249b.md)

## result

PASS

## statement

Verified f5cb4e13e at HEAD: re-ran its Checks: go test -trimpath -p 2 ./internal/run/..., ./internal/api/..., ./internal/guide/... exit 0; ./cmd/rhiza -run TestScripts/describe-golden -update exit 0 with no golden diff; -run 'TestScripts/(describe-golden|run-review-journal)' exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0. Diff: run.Kinds ends oversized-leaf, no-failure-case and KindSeverity makes it major (TestKindsPlanReviewAdditions); TestRunFindingsPlanKinds runs each plan kind in its own run's round 1 so a confirmed no-failure-case fails review-check before round 3, and the unknown kind no-failure is refused with exit 2 listing all eight kinds and journals nothing; TestSeveritiesBlockShared needs the severities block to name no-failure-case major; run-review-journal.txtar, describe-golden and docs/reference/cli.md list the eighth kind. Scope: the severities.md rewording keeps every severity and per-check mapping and is in the leaf's allowed files; the plan_reviewer_test.go skip of no-failure-case is needed because the leaf must not change plan-reviewer's text, and the plan-reviewer-leaf-contract leaf's acceptance makes that topic list eight kinds. No stubs or TODOs. Nit: the RunFinding kind desc lacks a comma after '(an intent clause no leaf serves)'.

## Links

### Out

- verifies [no-failure-case-kind](../capability/rhiza/workflows/built-in-workflows/plan-workflow/no-failure-case-kind.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
