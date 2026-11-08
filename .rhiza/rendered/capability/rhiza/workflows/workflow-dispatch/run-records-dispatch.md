# run-records-dispatch

capability · `01M6VYTB1Y3D4W01CFE712D6TD` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript shows run start on a change set picked by a [ready] rule recording rule ready.bug, one with no dispatch file recording the default, and one with --workflow recording explicit, each shown by run status --json and run ls --json; a test in internal/api shows the daemon's start log line naming the pick

## statement

Each run's journal records which workflow dispatch picked and why (the request's or change set's override, the rule trigger.kind, the built-in default, or an explicit --workflow), rhiza run ls and rhiza run status show it, and the daemon logs it when it starts the run

## Links

### Out

- depends-on [request-trigger-dispatch](request-trigger-dispatch.md)
- depends-on [ready-trigger-dispatch](ready-trigger-dispatch.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- verifies [run-records-dispatch-verified](../../../../evidence/run-records-dispatch-verified.md)
- verifies [run-records-dispatch-qa-1](../../../../evidence/run-records-dispatch-qa-1.md)
- verifies [run-records-dispatch-qa-2](../../../../evidence/run-records-dispatch-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
