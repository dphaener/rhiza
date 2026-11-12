# ship-change-model-conflict-fix-step-verified

evidence · `01M77QEJNYRTVS026M8W9Q5H7R` · actor [claude-code-19e76a99](../actor/claude-code-19e76a99.md)

## result

PASS

## statement

At 04706495d: go test -trimpath ./internal/workflow/... passed (TestShipChangeModelConflictFix asserts land-triage pass gate-fix-merge/fail land-conflict-triage/max_visits 4, land-conflict-triage running rhiza ship check --model-conflict with only on.pass model-conflict-fix and max_visits 4, model-conflict-fix an agent step passing only to gate-land with max_visits 3, land 4, gate-land 9, and the prompt naming rhiza set <node> --rebase and the Change/Actor/Review/Model-Conflict trailers). go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-model-conflict passed: with a bare origin, fake gh and fake fixer, land is refused for the model conflict, the fixer commits one rebased update, the run returns through gate-land (visit 2), lands and merges; a refused push with no model conflict fails at land-conflict-triage with no fixer session. rhiza workflow check passed on every built-in toml and on the repo's ship-change, plan-bug and bugfix (only the expected bypass warnings). Read the diff: the repo's .rhiza/workflows/ship-change.toml carries the same steps and caps as the built-in.

## Links

### Out

- verifies [ship-change-model-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-model-conflict-fix-step.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
