# dispatch-guide-topic

capability · `01M6VYTB5F2Z10QBJNJH3M0BGN` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A test of the guide shows the workflows topic naming dispatch.toml, [request], [ready] and request --workflow, and the guide's existing size checks still pass

## statement

The workflows guide topic explains .rhiza/dispatch.toml: its two tables, the kinds, the built-in defaults, rhiza request --workflow, the resolution order and how a broken name is reported

## Links

### Out

- depends-on [request-workflow-flag](request-workflow-flag.md)
- depends-on [ready-trigger-dispatch](ready-trigger-dispatch.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- verifies [dispatch-guide-topic-verified](../../../../evidence/dispatch-guide-topic-verified.md)
- verifies [dispatch-guide-topic-qa-1](../../../../evidence/dispatch-guide-topic-qa-1.md)
- verifies [dispatch-guide-topic-qa-2](../../../../evidence/dispatch-guide-topic-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
