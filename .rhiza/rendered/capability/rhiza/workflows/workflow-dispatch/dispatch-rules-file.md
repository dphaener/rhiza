# dispatch-rules-file

capability · `01M6VYTACHSRCG9FCR50S9S8JB` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

Unit tests in the package that parses the file cover: no file gives plan for every request kind and ship-change for every ready kind; a file setting only [ready] bug changes only that pick; a file holding [request] and [[schedules]] reads with no error; an unknown table, an unknown kind key and a non-string or path-like value each fail with an error naming the table and key

## statement

Rhiza reads the repo's committed .rhiza/dispatch.toml: a [request] table (a request is queued) and a [ready] table (a change set is accepted and ready), each mapping a kind (feature, bug, milestone) to a workflow name, beside the [[schedules]] entries schedules-file defines; a missing file, table or key falls back to the built-in default (request: plan; ready: ship-change), and any table other than [request], [ready] and [[schedules]], any other key in [request] or [ready], or a value that is not a bare workflow name, is an error naming it

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
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
