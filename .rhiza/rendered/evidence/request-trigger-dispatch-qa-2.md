# request-trigger-dispatch-qa-2

evidence · `01M6WB8WP595HE9V7XA41SWX4Z` · actor [claude-code-84def52f](../actor/claude-code-84def52f.md)

## result

PASS

## statement

QA cycle 2 at 76c40f51, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows plan/bugfix/other in .rhiza/workflows, RHIZA_* unset): With [request] bug = "bugfix" and milestone = "other", the daemon logs 'starting the plan run ... with workflow bugfix (picked by rule request.bug)' and detaches run start --workflow bugfix --pick rule request.bug --new ...; a feature request gets plan (picked by default); the milestone request gets --workflow other --milestone, without --autonomous. run start --foreground --workflow bugfix --new ... --request ran to success, and the request's change set records request_kind bug.

## Links

### Out

- verifies [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
