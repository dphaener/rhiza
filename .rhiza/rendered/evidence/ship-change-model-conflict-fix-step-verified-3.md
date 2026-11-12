# ship-change-model-conflict-fix-step-verified-3

evidence · `01M77RRCS1Y9785E5ZRXXEZ02V` · actor [claude-code-a507bd58](../actor/claude-code-a507bd58.md)

## result

PASS

## statement

Attempt 3 at 51867cf26 (fix for the model check's major finding): read the diff, which restates ship-change-prepush-fix (statement and acceptance now say land-triage goes on.fail only to land-conflict-triage, open-pr's and push-merge's triages have no on.fail) and links this leaf depends-on it, plus the actor records; checked built-in ship-change.toml: land-triage max_visits 4 on.pass gate-fix-merge on.fail land-conflict-triage, land-conflict-triage max_visits 4 on.pass model-conflict-fix with no on.fail, model-conflict-fix max_visits 3 on.pass gate-land, land max_visits 4, gate-land max_visits 9, open-pr-triage with no on.fail, so the restated node matches the code. Ran go test -trimpath ./internal/workflow/... (ok), go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-model-conflict' (PASS), go test -trimpath ./internal/api -run TestShipGate (ok), rhiza workflow check ship-change, bugfix and plan-bug (all valid), and make self-check in a clean detached worktree at HEAD: no fold or model errors, only stale .rhiza/rendered files (review sessions' actor pages and index.md) that the run renders at close-out.

## Links

### Out

- verifies [ship-change-model-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-model-conflict-fix-step.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
