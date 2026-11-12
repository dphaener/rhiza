# ship-change-model-conflict-fix-step

capability · `01M77NA4TTXHP9DFK9NK1Z1453` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/workflow loads the built-in ship-change workflow and asserts land-triage on.pass gate-fix-merge, on.fail land-conflict-triage and max_visits 4, land-conflict-triage running rhiza ship check --model-conflict with on.pass model-conflict-fix, no on.fail and max_visits 4, model-conflict-fix an agent step with on.pass gate-land and max_visits 3, land max_visits 4, gate-land max_visits 9, and the prompt naming rhiza set --rebase and the trailers; a testscript run-ship-model-conflict.txtar with a local bare origin, a fake gh and a fake fixer agent: another change set lands a write to the same node field before land, land is refused, the fake fixer commits a rebased update, the run goes back through gate-land and lands once and merges; a land failing on a refused push (no model conflict) ends the run at land-conflict-triage with no fixer session; rhiza workflow check passes on the built-ins. Run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-model-conflict' only.

## statement

The built-in ship-change workflow recovers from a land refused for a model conflict: land-triage (rhiza ship gate --failed) keeps on.pass gate-fix-merge and gains on.fail land-conflict-triage, a command step running rhiza ship check --model-conflict with on.pass model-conflict-fix and no on.fail, so a gh, git or network failure of land still ends the run; model-conflict-fix is an agent step on the run worktree whose prompt has the fixer read both change sets' intent from rhiza ship check --json, move each stale update onto the node's landed version with rhiza set <node> --rebase setting every field it writes so both sides' edits hold, run rhiza check, commit with Change:, Actor:, Review: <run> and Model-Conflict: trailers, never push, and on an edit that no statement can honour on both sides ask one blocking question naming both change sets and the node, then run rhiza session done; it passes to gate-land, so land pushes the fix and lands once. The caps fit 3 fixes followed by a land that fails its pre-push check: model-conflict-fix has max_visits 3, land-conflict-triage max_visits 4 (a 4th conflict reaches it and ends the run on ship check's '3 fixes made' reason, not on the cap), land-triage rises from 3 to 4, land from 1 to 4, and gate-land from 6 to 9.

## Links

### Out

- depends-on [ship-land-records-model-conflict](ship-land-records-model-conflict.md)
- depends-on [ship-change-prepush-fix](ship-change-prepush-fix.md)
- decided-by [ship-land-model-conflict-rebase](../../../../../decision/ship-land-model-conflict-rebase.md)

### In

- verifies [ship-change-model-conflict-fix-step-verified](../../../../../evidence/ship-change-model-conflict-fix-step-verified.md)
- verifies [ship-change-model-conflict-fix-step-qa-1](../../../../../evidence/ship-change-model-conflict-fix-step-qa-1.md)
- verifies [ship-change-model-conflict-fix-step-verified-3](../../../../../evidence/ship-change-model-conflict-fix-step-verified-3.md)
- verifies [pr-484-fix-watchdog-unmerged-pr-at-01m77n1gagbp](../../../../../evidence/pr-484-fix-watchdog-unmerged-pr-at-01m77n1gagbp.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
