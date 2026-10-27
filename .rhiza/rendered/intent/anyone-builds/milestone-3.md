# milestone-3

intent · `01M5WPTJ9Z4Z1FWZC11QMCPV2D` · parent [anyone-builds](../anyone-builds.md)

## milestone

3

## statement

Rhiza's own workflows: milestone 4 is planned, built and reviewed by Rhiza, with the owner talking only to Rhiza.

## Links

### Out

- decided-by [m3-rhizas-own-workflows](../../decision/m3-rhizas-own-workflows.md)
- decided-by [model-is-plan-and-state](../../decision/model-is-plan-and-state.md)
- decided-by [four-workflows](../../decision/four-workflows.md)
- decided-by [two-owner-gates](../../decision/two-owner-gates.md)
- decided-by [autonomous-mode](../../decision/autonomous-mode.md)
- decided-by [plan-workflow-shape](../../decision/plan-workflow-shape.md)
- decided-by [build-workflow-shape](../../decision/build-workflow-shape.md)
- decided-by [verifier-writes-evidence](../../decision/verifier-writes-evidence.md)
- decided-by [full-agent-sessions](../../decision/full-agent-sessions.md)
- decided-by [engine-is-conductor](../../decision/engine-is-conductor.md)
- decided-by [one-worktree-first](../../decision/one-worktree-first.md)
- decided-by [sessions-talk-freely](../../decision/sessions-talk-freely.md)
- decided-by [arbiter](../../decision/arbiter.md)
- decided-by [review-workflow-shape](../../decision/review-workflow-shape.md)
- decided-by [m3-scope](../../decision/m3-scope.md)
- decided-by [m3-change-sets](../../decision/m3-change-sets.md)
- decided-by [budget-timed-alone](../../decision/budget-timed-alone.md)
- derived-from [milestone-3-conversation](../../source/milestone-3-conversation.md)

### In

- serves change set `01M5WPTRX1DDQ7BAGDF7T29QPW` (Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.)
- serves change set `01M5WPTV0NX6YECXCPMHSEJAWY` (Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.)
- serves change set `01M5WPTVT93DTC076WJNB7E2KX` (Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.)
- serves change set `01M5WPTWM0MPSBEHRYDMV9Q7PM` (Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.)
- serves change set `01M5X3GBH4ZQTZVXFEHR966747` (ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.)

## History

- `01M5WPTHR3EBHDNNH1CFX7A0T0` Milestone 3 planning record: the conversation, the milestone-3 intent and its done points, and its decisions
