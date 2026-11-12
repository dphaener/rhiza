# ship-change-model-conflict-fix-step-qa-1

evidence · `01M77QZV2VRVZE4N2X6Q3YZ4K0` · actor [claude-code-b5a348d0](../actor/claude-code-b5a348d0.md)

## result

PASS

## statement

QA with a binary built from the branch head: rhiza workflow check ship-change reports the workflow valid (only the expected bypass warning, now naming model-conflict-fix); rhiza workflow describe ship-change shows land on.fail land-triage, land-triage on.fail land-conflict-triage, land-conflict-triage on.pass model-conflict-fix with no on.fail, and model-conflict-fix an agent step on.pass gate-land whose prompt names rhiza set --rebase and the Change:, Actor:, Review: and Model-Conflict: trailers; the repo's .rhiza/workflows copy matches the built-in from land to wait-landed. go test -trimpath ./internal/workflow/... and the run-ship-model-conflict and describe-golden testscripts pass (the fixer case goes land fail, land-triage fail, land-conflict-triage pass, one model-conflict-fix, gate-land, land pass and merges; a refused push ends at land-conflict-triage with no fixer).

## Links

### Out

- verifies [ship-change-model-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-model-conflict-fix-step.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
