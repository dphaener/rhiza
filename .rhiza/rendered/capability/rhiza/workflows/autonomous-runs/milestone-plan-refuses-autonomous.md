# milestone-plan-refuses-autonomous

capability · `01M6MYGJT14RGWC73S0V3A5YEB` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

Testscripts: run start --workflow plan --new <intent> --milestone --autonomous, and run start --autonomous on a record change set, each exit with a usage error naming plan-milestone's declaration and write no journal and no draft; the same starts without --autonomous run; run start --workflow plan --new <intent> --autonomous (an intention) still runs autonomous; with a repo plan-milestone.toml declaring autonomous = true the milestone start with --autonomous runs; run resume of an autonomous-journaled run whose workflow declares owner is refused

## statement

rhiza run start refuses --autonomous on a run whose workflow declares the stage it performs owner (by default the built-in plan-milestone, which milestone plan runs use) with a usage error naming the workflow and its declaration, and run resume refuses to continue such a run journaled autonomous; a repo that overrides plan-milestone to declare auto lifts the refusal

## Links

### Out

- decided-by [milestone-plan-never-autonomous](../../../../decision/milestone-plan-never-autonomous.md)
- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [milestone-planning-always-plan](../../../../source/milestone-planning-always-plan.md)

### In

- verifies [milestone-plan-refuses-autonomous-verified](../../../../evidence/milestone-plan-refuses-autonomous-verified.md)
- verifies [milestone-plan-refuses-autonomous-qa-1](../../../../evidence/milestone-plan-refuses-autonomous-qa-1.md)
- verifies [milestone-plan-refuses-autonomous-verified-2](../../../../evidence/milestone-plan-refuses-autonomous-verified-2.md)
- verifies [milestone-plan-refuses-autonomous-qa-2](../../../../evidence/milestone-plan-refuses-autonomous-qa-2.md)
- verifies [pr-289-milestone-planning-always](../../../../evidence/pr-289-milestone-planning-always.md)
- verifies [milestone-plan-refuses-autonomous-qa-3](../../../../evidence/milestone-plan-refuses-autonomous-qa-3.md)
- verifies [milestone-plan-refuses-autonomous-qa-2-01m6vy7h](../../../../evidence/milestone-plan-refuses-autonomous-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
