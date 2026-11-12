# ship-model-conflict-fix-verified

evidence · `01M77R83R5H326P046ZS6DG7CK` · actor [claude-code-53c2b18e](../actor/claude-code-53c2b18e.md)

## result

PASS

## statement

Verifier at 4f2a692c5: go test -trimpath ./internal/workflow/... ok, including the new TestShipChangeModelConflictRoute (land on.fail -> land-triage -> model-conflict-triage running `rhiza ship check --model-conflict`, pass -> model-conflict-fix agent whose prompt names `rhiza set <ref> --rebase` and the Review: {run} trailer -> render-model-conflict -> gate-land -> land; land max_visits=3, every new step bounded, triage fail ends the run). go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-model-conflict-fix$|TestScripts/run-ship-conflict-fix$' -count=1 ok: case 1 land refused on the stale-base update, fixer runs `rhiza set first --rebase` and commits with Review: trailer, gate-land and land visit 2 pass, PR merged (Ship it (#7)), first reads 'Theirs and ours'; case 2 a fixer that repairs nothing fails at model-conflict-fix with no merge. The repo's .rhiza/workflows/ship-change.toml carries the same new steps (diff vs builtin is only the pre-existing head-build steps). make self-check: the model folds clean; in a detached worktree at the leaf commit, after `rhiza render` (the run's render step), the only fail is docs/reference/cli.md, which fails identically at the base a4fd2b7 (pre-existing, not this leaf). Diff read: workflow TOML, test updates to checkShipChangePrePushFix and TestShipChangeLocalGate, new txtar.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
