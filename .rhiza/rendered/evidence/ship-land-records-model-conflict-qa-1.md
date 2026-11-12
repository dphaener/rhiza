# ship-land-records-model-conflict-qa-1

evidence · `01M77QZV0JMEF9DC302ABJ83SS` · actor [claude-code-b5a348d0](../actor/claude-code-b5a348d0.md)

## result

PASS

## statement

QA with a binary built from the branch head, in a scratch repo with a local bare origin and a fake gh: change set Ship it updated store's statement, then change set Other landed its own statement write. rhiza land --json exited 3 with error.conflict naming node, op, other change set and aspects [field:statement]; the human message is unchanged. rhiza ship check then rhiza ship land exited 3 with the same body and the ship state gained model_conflict (node, op, change_set, aspects, head). ship check --model-conflict --json passed (exit 0) with refused_change_set, intent 'Ship it', change_set and other_intent 'Other'; it failed with 'no model conflict is recorded' before any land, failed with 3 commits carrying Review: RUNX and Model-Conflict: when RHIZA_RUN=RUNX (passed for another run ID), failed once the branch moved past the recorded head, and refused --conflicts beside it (exit 2). After rhiza set store --rebase and a trailer commit, ship land landed, pushed and cleared model_conflict from the state. Scoped tests internal/api -run 'Ship|Land' (from a worktree outside ~/code) and internal/model pass. Minor: the passing human line reads 'open from cs/ship-it', not the conflict.

## Links

### Out

- verifies [ship-land-records-model-conflict](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-land-records-model-conflict.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
