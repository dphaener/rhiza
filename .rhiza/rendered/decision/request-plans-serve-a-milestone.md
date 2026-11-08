# request-plans-serve-a-milestone

decision · `01M6JDXH6A7W60623XN96ZV9JV`

## reasoning

The change set's own intention says a request's plan should link to 'the open milestone it fits, or a backlog milestone the request names', and milestone-25's statement is exactly that every change set planned from queued requests is built and shipped by the daemon, so this fix serves it, which also lets the daemon build it. RequestRequest (internal/api/request.go) today has only text and kind, so naming a milestone needs a new optional flag; leaving it optional keeps plan=auto with accept=auto working without owner input. The daemon already skips with daemon.SkipNoMilestone (internal/daemon/pick.go, internal/api/daemon.go), but only reports it in status after the fact; the failure that prompted this (paging 01M6H3A65RDP, usage governor 01M6HHMTEDFK) was silent, so the check belongs where the plan is made, failing the plan step before acceptance. The owner files requests from the Rex TUI apps (owner goal: run Rhiza from TUI apps), so the requests app needs the same optional milestone field the CLI gets; defaulting it to auto keeps filing a one-line act.

## statement

Change set 01M6JDM serves milestone-25 (01M6HKVRE1B7). A request's plan reaches a milestone in two ways: `rhiza request` (and the request op) gains an optional --milestone flag naming a milestone intent by handle or ID, stored on the request and handed to the planner; with none named, the planner links the change set serves to the open milestone it fits (a milestone intent, never an umbrella or area intent under no milestone). The guard is a command step right after the plan's interview: it fails the run, loudly, when a planned change set serves no milestone directly or through an intent under one, so a plan can no longer be accepted and then silently skipped. (Not the interview's until check: the engine only fails an agent step on a failed until once its session exits, so a live session that said done would just be polled again.) The requests TUI app (q) gets an optional milestone field when filing a request, defaulting to auto (planner picks the open milestone it fits).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-6b593167-1](../question/question-6b593167-1.md)
- derived-from [plan-request-plans-serve-milestone](../source/plan-request-plans-serve-milestone.md)

### In

- decided-by [request-names-milestone](../capability/rhiza/orchestration/daemon/plan-requests/request-names-milestone.md)
- decided-by [plan-run-hands-milestone](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-hands-milestone.md)
- decided-by requests-app-milestone (retired)
- supersedes [scheduling-policy-shape](scheduling-policy-shape.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
