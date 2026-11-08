# ready-trigger-dispatch

capability · `01M6VYTAW1NNHC09NTGF6RPVKS` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

Tests in internal/api show the resolver's order (recorded workflow over rule over ship-change, explicit --workflow over all), and a testscript shows run start --dry-run on a bug change set with [ready] bug set to a repo workflow naming that workflow, and with no dispatch file naming ship-change

## statement

One resolver picks the workflow of a ready change set's run, used both by the daemon's ship pick and by rhiza run start <cs> with no --workflow: the change set's recorded workflow, else the [ready] rule for its kind, else ship-change; an explicit run start --workflow always wins

## Links

### Out

- depends-on [dispatch-rules-file](dispatch-rules-file.md)
- depends-on [change-set-records-dispatch](change-set-records-dispatch.md)
- decided-by [workflow-dispatch-shape](../../../../decision/workflow-dispatch-shape.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- depends-on [dispatch-unresolved-attention](dispatch-unresolved-attention.md)
- depends-on [run-records-dispatch](run-records-dispatch.md)
- depends-on [dispatch-guide-topic](dispatch-guide-topic.md)
- verifies [ready-trigger-dispatch-verified](../../../../evidence/ready-trigger-dispatch-verified.md)
- verifies [ready-trigger-dispatch-qa-1](../../../../evidence/ready-trigger-dispatch-qa-1.md)
- verifies [ready-trigger-dispatch-verified-2](../../../../evidence/ready-trigger-dispatch-verified-2.md)
- verifies [ready-trigger-dispatch-qa-2](../../../../evidence/ready-trigger-dispatch-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
