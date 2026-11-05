# plan-gate

capability · `01M5WSGDEE854PFE4BXJM13CWS` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript ends a plan run whose change set adds an intent served by two draft change sets: normal (a record change set, landed after the interview) with rhiza.autonomy.accept owner, both are proposed, none accepted, the landed record is left as it is, and run status --json names rhiza accept --rule owner; normal with accept auto and a passing review, both are accepted with rule autonomy; autonomous (an intention change set, since a milestone plan never runs autonomous) with a passing review, both are accepted with rule autonomous and reasons naming the run; autonomous with a question left by the review, none is accepted and the command exits non-zero naming the question; once the arbiter's Decision answers it, both are accepted. A run on an intention change set that serves an existing intent proposes that change set itself.

## statement

rhiza run gate ends a plan: it proposes each draft among the plan's change sets (the run's change set and every change set that serves an intent the run's change set adds), leaving landed or accepted ones as they are. In a normal run it follows rhiza.autonomy.accept: when it is auto and the plan's review passed, it accepts each with rule autonomy; when it is owner, or the review has not passed, they wait for the owner's rhiza accept --rule owner, which run status names. In an autonomous run it accepts each with rule autonomous once the review passed and every question from the review is answered, and accepts nothing while a question from the review is open.

## Links

### Out

- depends-on [plan-review-rounds](plan-review-rounds.md)
- depends-on [autonomous-gate-rule](../../autonomous-runs/autonomous-gate-rule.md)
- decided-by [two-owner-gates](../../../../../decision/two-owner-gates.md)
- decided-by [autonomous-mode](../../../../../decision/autonomous-mode.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [plan-review-left-arbitrated](../../../../../decision/plan-review-left-arbitrated.md)

### In

- depends-on [plan-workflow-file](plan-workflow-file.md)
- relates-to [plan-record-landing](../../../../../question/plan-record-landing.md) — The gate proposes the plan's own change set with the others; who lands it after acceptance is not decided.
- relates-to [plan-review-left-autonomous](../../../../../question/plan-review-left-autonomous.md) — The gate accepts nothing while a review question is open; whether the arbiter settles those questions first is not decided.
- verifies [plan-gate-verified](../../../../../evidence/plan-gate-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)
- verifies [pr-289-milestone-planning-always](../../../../../evidence/pr-289-milestone-planning-always.md)
- verifies [pr-315-fix-an-owner-run-plan-s-gate](../../../../../evidence/pr-315-fix-an-owner-run-plan-s-gate.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
