# dispatch-rules-check

capability · `01M6VYTAG5STEEVHYYJ6V971SN` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript (dispatch-check) shows rhiza check passing with no dispatch file, with a rule naming a repo workflow in .rhiza/workflows/ and with a rule naming a built-in workflow, and failing with a rule naming a missing workflow, with a rule naming a workflow that exists only in the personal layer ($XDG_CONFIG_HOME/rhiza/workflows set inside the script), and with an unknown table, each failure naming trigger.kind or the table

## statement

rhiza check fails on a .rhiza/dispatch.toml that does not parse, or whose rule names a workflow that neither the repo layer (.rhiza/workflows/) nor the built-in layer resolves, naming the file, the trigger.kind and the workflow; it never reads the personal layer ($XDG_CONFIG_HOME/rhiza/workflows), so the same committed .rhiza/ gives the same result on every machine and in CI

## Links

### Out

- depends-on [dispatch-rules-file](dispatch-rules-file.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- verifies [dispatch-rules-check-verified](../../../../evidence/dispatch-rules-check-verified.md)
- verifies [dispatch-rules-check-qa-1](../../../../evidence/dispatch-rules-check-qa-1.md)
- verifies [dispatch-rules-check-verified-2](../../../../evidence/dispatch-rules-check-verified-2.md)
- verifies [dispatch-rules-check-qa-2](../../../../evidence/dispatch-rules-check-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
