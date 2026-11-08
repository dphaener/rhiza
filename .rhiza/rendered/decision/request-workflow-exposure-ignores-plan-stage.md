# request-workflow-exposure-ignores-plan-stage

decision · `01M6W63782G9RRJWP7E5AQAPJK`

## reasoning

Read the question, the leaf request-exposed-workflows, decision risk-warnings-shape and the code. internal/api/request.go daemonPlanRequests starts a plan run for every queued request with no check of the plan stage (planMode only decides --autonomous), and planRunArgs always passes --permission-mode <rhiza.daemon.permissionMode> and --ref; even milestone requests, which are never autonomous, start the same way. internal/workflow/builtin/plan.toml's interview step starts an agent at once whose prompt carries {title}, the request text. So the decision's premise that with plan owner 'a person reads the plan before anything runs' is false for the plan run itself: the agent acts first, in bypass mode, on outside text. The change set's intent is to warn when outside text reaches agents with no permission checks, and it never blocks, so a false negative is worse than an extra warning; the warning stays warn and only reads config. Justifying the interactive run as gated would need the agent to wait for the owner before any tool call, which nothing enforces. The ready workflow is different: under plan owner its input is a plan a person shaped, and under accept owner a person accepts it first, so its stricter rule still holds.

## statement

The request workflow (plan today, or whatever a [request] rule names) is flagged whenever any agent session it starts runs in a bypass mode, whatever rhiza.autonomy.plan says, and milestone requests count too. An interactive plan run is not a gate: the daemon starts it from every queued request at once, its interview agent's prompt carries the request text and its change set the --ref, and under a bypass mode that agent can act on them before the owner joins; the plan workflow's review, critic and reviser agents read them too. This replaces the 'when rhiza.autonomy.plan is auto' condition in @risk-warnings-shape point (2) and @request-exposed-workflows for the request workflow, and the 'plan owner flags nothing' acceptance case becomes 'plan owner with a bypass daemon mode flags plan, not ship-change'. The ready workflow's rule is unchanged: it is flagged only when rhiza.autonomy.plan, accept and ship are all auto, because with plan owner a person co-writes the plan, and with accept owner a person accepts it, before the ready workflow's agents see anything.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6w4k1-r3-f1](../question/review-01m6w4k1-r3-f1.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
