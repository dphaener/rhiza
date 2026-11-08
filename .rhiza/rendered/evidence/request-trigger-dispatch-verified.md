# request-trigger-dispatch-verified

evidence · `01M6W2XJGWA5D60ERS91EA263P` · actor [claude-code-f1847a18](../actor/claude-code-f1847a18.md)

## result

PASS

## statement

At 1a855d1c I read the diff of 8c465daa and 1a855d1c and ran go test -trimpath -count=1 ./internal/api -run 'OtherWorkflow|DispatchedWorkflow'. All pass. TestDaemonStartsRequestWithDispatchedWorkflow: with a dispatch.toml of [request] bug=bugfix and milestone=roadmap, a bug request starts --workflow bugfix --new, a feature request starts --workflow plan, and a milestone request starts --workflow roadmap --milestone without --autonomous. TestRunResumeGatesMilestonePlanUnderOtherWorkflow: an autonomous Plan run under roadmap on a record change set hits the milestone owner gate. TestInboxKeepsPlanRunUnderOtherWorkflow: the inbox keeps it listed (not skipped). TestRunResumeLetsLandedPlanRunUnderOtherWorkflow: run.resume lets a Plan run under bugfix on a landed change set through, and refuses a non-plan run that never reached land. This answers attempt 1's finding f1. TestResumablePlanRunUnderOtherWorkflow in internal/daemon shows Resumable keeps the landed plan run under bugfix. ./internal/daemon, ./internal/run, ./internal/dispatch and the cmd/rhiza scripts run-start-new, request* and daemon* all pass. The keying is run.State.IsPlanRun (Started.Plan, set when the run is started with --new, or a legacy plan workflow). It is used in Resumable, inboxItems and both run.resume checks.

## Links

### Out

- verifies [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
