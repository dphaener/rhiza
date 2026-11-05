# milestone-plan-never-autonomous

decision · `01M6MYD57G39XAW6AF3CXS8XA8`

## reasoning

milestone-22's statement is the one about each stage (planning included) switching between autonomous and owner-driven, and this change set narrows how the plan stage's autonomy applies; milestone-25 is about emptying the backlog, which this does not serve. In the code, the daemon decides autonomy in daemonPlanRequests/daemonStartPlan (internal/api/request.go) from autonomyStage(plan) alone, adding --milestone for KindMilestone requests; run start (internal/api/run.go runRunStart/newPlanChange) maps --milestone to a record change set and accepts --autonomous with it. Gating only the daemon would leave a hand-started or resumed run able to go autonomous, so the owner's rule ('that always requires a human') needs the run start refusal too; a refusal there is explicit, while the daemon simply drops the flag so its milestone requests still start and wait for the owner.

## statement

The change set serves milestone-22 (switching each stage between autonomous and owner-driven). Rhiza enforces it in two places: (1) the daemon's daemonStartPlan never passes --autonomous for a request of kind milestone, whatever rhiza.autonomy.plan says, so that plan run waits on the owner's interview; (2) run start refuses --autonomous on a plan run that plans a milestone (--milestone, a request of kind milestone, or an existing record change set) with a usage error naming the rule, so no other path, by hand or on resume, starts one autonomously. Feature and bug requests keep following rhiza.autonomy.plan.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-58fb9cbb-1](../question/question-58fb9cbb-1.md)

### In

- decided-by [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)
- decided-by [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
