# workflow-dispatch

capability · `01M6VYTAAZ3666PP9D4N9HTVSQ` · parent [workflows](../workflows.md)

## statement

Every run the daemon or run start begins picks its workflow through the repo's dispatch rules, keyed by trigger and request kind, with built-in defaults that reproduce plan and ship-change

## Links

### Out

- decided-by [workflow-dispatch-shape](../../../decision/workflow-dispatch-shape.md)

### In

- verifies [pr-367-workflow-dispatch](../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
