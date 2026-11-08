# request-workflow-flag-qa-1

evidence · `01M6W8JKDNP3NEN3HYVC31DHJ8` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): rhiza request 'Crash on empty input' --kind bug --workflow bugfix queued it; request ls --json shows "workflow":"bugfix" and the text output shows '[workflow bugfix]'; a plain request has no workflow field. request --workflow nothing-here exits 2 ('there is no such workflow in .rhiza/workflows, your personal workflows or built in') and queues nothing.

## Links

### Out

- verifies [request-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/request-workflow-flag.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
