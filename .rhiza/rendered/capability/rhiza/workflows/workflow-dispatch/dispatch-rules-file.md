# dispatch-rules-file

capability · `01M6VYTACHSRCG9FCR50S9S8JB` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

Unit tests in the package that parses the file cover: no file gives plan for every request kind and ship-change for every ready kind; a file setting only [ready] bug changes only that pick; an unknown table, an unknown kind key and a non-string or path-like value each fail with an error naming the table and key

## statement

Rhiza reads the repo's committed .rhiza/dispatch.toml: a [request] table (a request is queued) and a [ready] table (a change set is accepted and ready), each mapping a kind (feature, bug, milestone) to a workflow name; a missing file, table or key falls back to the built-in default (request: plan; ready: ship-change), and any other table or key, or a value that is not a bare workflow name, is an error naming it

## Links

### Out

- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)
- decided-by [workflow-dispatch-shape](../../../../decision/workflow-dispatch-shape.md)

### In

- depends-on [dispatch-rules-check](dispatch-rules-check.md)
- depends-on [request-trigger-dispatch](request-trigger-dispatch.md)
- depends-on [ready-trigger-dispatch](ready-trigger-dispatch.md)
- verifies [dispatch-rules-file-verified](../../../../evidence/dispatch-rules-file-verified.md)
- verifies [dispatch-rules-file-qa-1](../../../../evidence/dispatch-rules-file-qa-1.md)
- verifies [dispatch-rules-file-verified-2](../../../../evidence/dispatch-rules-file-verified-2.md)
- verifies [dispatch-rules-file-qa-2](../../../../evidence/dispatch-rules-file-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
