# milestone-plan-refuses-autonomous

capability · `01M6MYGJT14RGWC73S0V3A5YEB` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

Testscripts: run start --workflow plan --new <intent> --milestone --autonomous, and run start --workflow plan --autonomous on a record change set, each exit with a usage error naming the rule and write no journal and no draft; the same starts without --autonomous run, and run start --workflow plan --new <intent> --autonomous (an intention) still runs autonomous; run resume of a plan run on a record change set whose started entry is autonomous exits with the same usage error. plan-workflow-autonomous.txtar and plan-workflow-autonomous-left.txtar, which today start an autonomous plan run on a record change set, instead start it on an intention change set (its planner adds an intent without milestone=1, served by the draft it writes, and the -left script's landed-record Decision is one landed before the run), and still check what they check now: the arbiter answering the planner's question, review rounds, the check's questions answered where they live, and the gate accepting with rule autonomous; so plan-autonomous-interview, plan-review-arbitrated and plan-gate stay verified by them

## statement

rhiza run start refuses --autonomous on a plan run that plans a milestone (with --milestone, for a request of kind milestone, or on an existing record change set) with a usage error naming the rule that milestone planning needs the owner, and run resume refuses to continue such a run journaled autonomous, so no milestone plan run's agent gets RHIZA_AUTONOMOUS=1

## Links

### Out

- decided-by [milestone-plan-never-autonomous](../../../../decision/milestone-plan-never-autonomous.md)
- derived-from [milestone-planning-always-plan](../../../../source/milestone-planning-always-plan.md)

### In

- verifies [milestone-plan-refuses-autonomous-verified](../../../../evidence/milestone-plan-refuses-autonomous-verified.md)
- verifies [milestone-plan-refuses-autonomous-qa-1](../../../../evidence/milestone-plan-refuses-autonomous-qa-1.md)
- verifies [milestone-plan-refuses-autonomous-verified-2](../../../../evidence/milestone-plan-refuses-autonomous-verified-2.md)
- verifies [milestone-plan-refuses-autonomous-qa-2](../../../../evidence/milestone-plan-refuses-autonomous-qa-2.md)
- verifies [pr-289-milestone-planning-always](../../../../evidence/pr-289-milestone-planning-always.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
