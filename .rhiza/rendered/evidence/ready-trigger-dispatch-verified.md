# ready-trigger-dispatch-verified

evidence · `01M6W5D1T25YYFM8VMZ4TTPZVY` · actor [claude-code-5b866a73](../actor/claude-code-5b866a73.md)

## result

PASS

## statement

At 044ec30f I read the diff df9cc6f4..HEAD: one resolver readyWorkflow (recorded workflow, else the [ready] rule for its kind, else ship-change) used by the daemon's ship pick (daemonShipOf For, daemonStart passes --workflow) and by run start on an existing change set via startWorkflow (explicit --workflow over all). go test -trimpath ./internal/api -run 'TestReadyWorkflowOrder|TestStartWorkflowExplicitWins|TestRunStartExplicitWorkflowOverRule' passes: recorded over rule over ship-change, and explicit --workflow over a recorded workflow, a rule and the default. go test -trimpath ./internal/daemon/... passes (TestPickShippedUnderDispatchedWorkflow). go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start-ready-dispatch|describe-golden|request-workflow-flag)$' passes: run-start-ready-dispatch makes a real bug change set (request --kind bug, plan run), and run start --dry-run with no dispatch file fails naming ship-change ('ship-change builds and reviews its leaves'), with [ready] bug = "bugfix" prints '- Workflow: bugfix', and with --workflow other prints '- Workflow: other'. Attempt 1's three findings are fixed.

## Links

### Out

- verifies [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
