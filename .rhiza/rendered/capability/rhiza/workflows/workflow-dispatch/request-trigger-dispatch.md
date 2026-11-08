# request-trigger-dispatch

capability · `01M6VYTARWAQRW1FDT2AKAFZQA` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

Tests in internal/api show a dispatch file mapping [request] bug to a repo workflow making the daemon's start of a bug request run that workflow with --new, a feature request still running plan, and a milestone request mapped to another workflow starting it with --milestone and without --autonomous; tests in internal/daemon and internal/api show a request's run under a workflow not named plan still resumed by daemon.Resumable after its change set lands early, still listed in the inbox, and a milestone request's run under such a workflow still stopped by run resume's owner gate

## statement

The daemon starts a queued request's plan run with the workflow the [request] rule names for the request's kind (default plan), and run start --new accepts that workflow instead of refusing every workflow but plan; a milestone request still opens a record change set and never runs autonomous, whatever workflow its rule names. Whether a run is a request's plan run is keyed on the run having been started with --new or for a request, not on its workflow being named plan: daemon.Resumable (resuming a run whose change set landed early), the inbox's keeping such runs listed, and run resume's milestone owner gate all use that test

## Links

### Out

- depends-on [dispatch-rules-file](dispatch-rules-file.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- depends-on [dispatch-unresolved-attention](dispatch-unresolved-attention.md)
- depends-on [run-records-dispatch](run-records-dispatch.md)
- blocks [review-01m6vy7e-r3-f2](../../../../question/review-01m6vy7e-r3-f2.md)
- verifies [request-trigger-dispatch-verified](../../../../evidence/request-trigger-dispatch-verified.md)
- verifies [request-trigger-dispatch-qa-1](../../../../evidence/request-trigger-dispatch-qa-1.md)
- verifies [request-trigger-dispatch-qa-2](../../../../evidence/request-trigger-dispatch-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
