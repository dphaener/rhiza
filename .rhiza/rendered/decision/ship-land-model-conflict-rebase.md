# ship-land-model-conflict-rebase

decision · `01M77N8C8TWSDEB6TNWS8J8KD9`

## reasoning

Today land's on.fail goes to land-triage, which runs `rhiza ship gate --failed` and only passes a recorded pre-push failure, so a model conflict ends the run and needs an orchestrator by hand. The repair commit 480b6acf3 shows that this is routine and an agent can do it: one `rhiza set --rebase` that keeps both edits (trimGoCache removed, selfUpdate -> autoUpdate), acceptance unchanged. ship-change already handles the code equivalent this way: ship check --conflicts routes to conflict-fix, which merges both sides and asks a blocking question only for a real clash. A model conflict is the model-side version of that, and `rhiza set --rebase` exists for it and refuses unless every written field is set, so it cannot silently drop a side. Because land refuses before it lands anything, retrying through gate-land and land is safe and lands nothing twice. Ending the run means a person has to step in for a conflict that concurrent accepted change sets will keep producing, which goes against the goal that the daemon does all orchestration.

## statement

ship-change recovers from a land refused for a model conflict instead of ending the run. land-triage tells a model conflict (rhiza land exit 3, a stale base because another change set landed a write to the same node and field) apart from a gate failure and from gh/git/network errors, and routes it to a model-conflict fixer agent. The fixer reads both change sets' intent, moves each stale update onto the node's landed version with `rhiza set <node> --rebase`, setting every field so both sides' edits hold, commits it with the Change: and Actor: trailers, and the run goes back through gate-land to land. The fixer has a max_visits cap like the other fixers; if both edits genuinely clash it asks one blocking question naming both change sets and the node, as conflict-fix does for code. Every other land failure ends the run as now.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a86fdccf-1](../question/question-a86fdccf-1.md)
- derived-from [plan-ship-land-model-conflict](../source/plan-ship-land-model-conflict.md)

### In

- decided-by [ship-land-records-model-conflict](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-land-records-model-conflict.md)
- decided-by [ship-change-model-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-model-conflict-fix-step.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
