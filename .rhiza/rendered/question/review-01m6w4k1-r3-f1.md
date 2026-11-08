# review-01m6w4k1-r3-f1

question · `01M6W61C8JAH7D9G6J8G7HN3Z4`

## statement

Plan review round 3 left finding f1 (contradiction on request-exposed-workflows): The leaf, and Decision risk-warnings-shape point (2) it follows, flag the request workflow (plan) only when rhiza.autonomy.plan is auto, on the premise that with plan owner 'a person reads the plan before anything runs'. The code contradicts that. daemonPlanRequests (internal/api/request.go:308-336) starts the plan run from every queued request whatever the plan stage is. planRunArgs (request.go:387-402) always passes --permission-mode <rhiza.daemon.permissionMode> and the request's --ref, and only adds or omits --autonomous. So with plan owner the interview agent still starts at once in bypass mode, with the request text as its title and the ref in the change set, before any person has read it. The plan workflow's review, critic and reviser agents then read it too. As planned, a repo with plan owner and permissionMode bypassPermissions, which is exactly the outside-text exposure, gets no warning. The exposure test should not depend on rhiza.autonomy.plan for the request workflow, or the decision should say why an interactive plan run counts as gated even though its agent acts before the owner joins. The critic confirmed it: internal/api/request.go:336 starts a plan run for every queued request whatever rhiza.autonomy.plan is, and planRunArgs (request.go:387-402) always passes --permission-mode <rhiza.daemon.permissionMode> and --ref; plan owner only omits --autonomous. The plan workflow's interview step (internal/workflow/builtin/plan.toml) starts an agent at once whose prompt carries {title}, the request text, so under a bypass mode that agent acts on outside text before any person joins. The leaf's 'plan owner flags nothing' case and decision risk-warnings-shape point 2's premise are wrong: the request workflow's exposure must not hinge on rhiza.autonomy.plan alone (or the decision must justify why an interactive plan run counts as gated). How should it be settled?

## Links

### Out

- blocks [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)

### In

- answers [request-workflow-exposure-ignores-plan-stage](../decision/request-workflow-exposure-ignores-plan-stage.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
