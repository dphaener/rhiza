# run-records-dispatch-verified

evidence · `01M6W7ZDCASYR8CT722A6E6KJX` · actor [claude-code-2fea6e76](../actor/claude-code-2fea6e76.md)

## result

PASS

## statement

At fe6bb167 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-records-dispatch|describe-golden)$' (ok). The script run-records-dispatch.txtar starts a bug change set picked by the [ready] rule (run start, run status --json and run ls --json show pick rule ready.bug), a feature change set with no dispatch file (run start and run status --json show pick default, run ls --json lists it), and an explicit --workflow other run whose id is captured and checked by run status <id> --json (workflow other, pick explicit) and run ls --json, which fixes attempt 1's finding f1. In internal/api, TestDaemonStartsRequestWithDispatchedWorkflow asserts the daemon start log line 'with workflow bugfix (picked by rule request.bug)' and '(picked by default)' (daemon.go:1765); it, TestStartWorkflowExplicitWins and TestRunStartQueuesForTheDaemon pass; internal/dispatch and internal/run pass.

## Links

### Out

- verifies [run-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/run-records-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
