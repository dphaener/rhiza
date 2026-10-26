# ship-change-holds-up

decision · `01M5X3GBRZX8CEV9BHB88P666D`

## statement

The second proof run (01M5WXM5, change set 2) failed and showed four gaps in ship-change, fixed together before more change sets are run: the build step passed while shipit was still reviewing (idle between turns with a PR open); land rendered with the installed rhiza though the branch changed rendering; a PR that main moves past can conflict and get no CI; and a failed run can't be resumed.

## wording

the orchestrator's ruling as arbiter, under the owner's instruction to run milestone 3 and be the arbiter; not the owner's words

## Links

### Out

- derived-from [milestone-3-conversation](../source/milestone-3-conversation.md)

### In

- decided-by [agent-step-done-signal](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal.md)
- decided-by [land-with-branch-build](../capability/rhiza/workflows/workflow-runs/land-with-branch-build.md)
- decided-by [rebase-when-behind](../capability/rhiza/workflows/workflow-runs/rebase-when-behind.md)
- decided-by [resume-failed-run](../capability/rhiza/workflows/workflow-runs/resume-failed-run.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.
