# autonomy-per-workflow-scope

decision · `01M6VYJHQNTEEG7ERHY6X3E3D8`

## reasoning

milestone-25 is the backlog milestone the orchestrator opened under the owner's instruction of 2026-11-03 to empty the backlog and act on their behalf (decided on change set 01M6HKTEKH). This change set came in as a queued request that follows from decision workflows-are-the-rules, so it is backlog work, and milestone-25 is still in the live tree: planMilestoneOf (internal/api/plan.go:512) finds it, and milestone-check passed on retry once the link was added. rhiza health's roadmap count leaves it out, but that does not change what milestone-check reads. The arbiter's first ruling, to wait on 01M6VY7GRV, was written before this link existed; the link replaces that wait, and a landed milestone such as m26-autonomy-app is still not linked. An arbiter cannot open a milestone: milestone planning is owner-only (failed request 01M6MY89 quotes the owner: 'milestone planning should never be autonomous even in autonomous mode. that always requires a human'), and this change set's own intent keeps that as the milestone-planning default. Of the four follow-ups to decision workflows-are-the-rules, scheduling-as-policy (01M6VY7GRV) is the one that makes a milestone optional for a request's plan, so without milestone-25 the honest order would have been that one first, not a made-up milestone. For stages: today the plan run covers plan and accept at its gate, and a later ship-change run covers ship and signoff. Dispatch (01M6VY7G6H) records the chosen workflow on the change set so the ship pick follows it. Reading the declaration of whichever workflow runs a stage keeps each run self-contained, and it meets the acceptance: a bugfix workflow declaring autonomous = true goes through plan and ship with no owner step while rhiza.autonomy.plan stays owner. For milestone planning: a workflow file cannot declare autonomy per request kind, and the intent asks for a built-in milestone planning workflow that declares owner as a changeable default. So it has to be a separate workflow, and the --milestone hard-code has to go, or the default could not be changed.

## statement

(1) Change set 01M6VY7H serves milestone-25 (the backlog milestone: every change set planned from the owner's queued requests and the backlog is built and shipped by the daemon on its own). No new milestone is invented for it. The arbiter first ruled it would serve no milestone and wait on sibling 01M6VY7GRV (scheduling as policy); the orchestrator then linked it to milestone-25 (serves link 01M6W5PDXDBBEQZY4HHTDSSNH9), which is live, so milestone-check passes without waiting. (2) Each stage takes its autonomy from the workflow of the run that performs it: plan and accept from the plan run's workflow (the one dispatch picked for the request), ship and signoff from the ship run's workflow (the one picked for the ready change set, following the workflow the change set records). A declaration for a stage the run does not perform is ignored for that run, and a stage with no declaration falls back to rhiza.autonomy.<stage>. (3) Milestone planning gets its own built-in workflow, plan-milestone, with the plan workflow's steps, declaring autonomous = false. Milestone requests run it, through dispatch's default rule once 01M6VY7G6H lands and through daemonStartPlan until then. The hard-coded autonomous = false in daemonStartPlan (internal/api/request.go:366-371) goes, so a repo changes the default by overriding plan-milestone like any other workflow.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-e4df472e-1](../question/question-e4df472e-1.md)
- derived-from [autonomy-per-workflow-interview](../source/autonomy-per-workflow-interview.md)

### In

- decided-by [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)
- decided-by [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)
- decided-by [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)
- decided-by [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)
- decided-by [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)
- decided-by [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)
- decided-by [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
