# ship-model-conflict-fix-verified-2

evidence · `01M77TQ4TGSSCAC0ST0ESZXP7R` · actor [claude-code-21901b9f](../actor/claude-code-21901b9f.md)

## result

PASS

## statement

Verifier at 8b74e581c (attempt 2, after compliance f1): read the diff, which adds case 3 to run-ship-model-conflict-fix.txtar. In it the other change set lands only at gate-land through a fake local gate, so land is refused on the stale base. go test -trimpath ./internal/workflow/... passes, including TestShipChangeModelConflictRoute (land-triage to model-conflict-triage to model-conflict-fix, bounded, back through render-model-conflict). go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/run-ship-model-conflict-fix passes 3 times, plus run-ship-conflict-fix once. Case 3 asserts this route: gate-land v1 pass, land v1 fail, land-triage, model-conflict-triage pass, model-conflict-fix (rhiza set first --rebase, one fixer session), render-model-conflict, gate-land v2, land v2 pass. It also asserts the merged PR 'Ship it (#7)' and 'Theirs and ours' on first. Case 2 asserts that a fixer leaving the conflict fails the run unmerged. make self-check in a clean detached worktree at the commit fails only on 8 stale .rhiza/rendered files from the run's own actor records. After rhiza render, it passes.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
