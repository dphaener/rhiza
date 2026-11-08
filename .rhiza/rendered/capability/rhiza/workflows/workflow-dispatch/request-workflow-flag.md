# request-workflow-flag

capability · `01M6VYTAMKNBD6G38RFW88RFB7` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript shows request --workflow naming a repo workflow queued and listed by request ls --json with the workflow, request --workflow naming a missing workflow exiting 2 with the queue unchanged, and describe-golden regenerated for the new request field

## statement

rhiza request --workflow <name> stores the name on the queued request, refuses with a usage error (exit 2, nothing queued) a name no layer resolves, and request ls shows it; the request's plan run passes it to the change set it opens, so it overrides only that change set's ready pick

## Links

### Out

- depends-on [change-set-records-dispatch](change-set-records-dispatch.md)
- decided-by [workflow-dispatch-shape](../../../../decision/workflow-dispatch-shape.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- depends-on [dispatch-guide-topic](dispatch-guide-topic.md)
- verifies [request-workflow-flag-verified](../../../../evidence/request-workflow-flag-verified.md)
- verifies [request-workflow-flag-qa-1](../../../../evidence/request-workflow-flag-qa-1.md)
- verifies [request-workflow-flag-qa-2](../../../../evidence/request-workflow-flag-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
