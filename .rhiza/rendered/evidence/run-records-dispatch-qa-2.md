# run-records-dispatch-qa-2

evidence · `01M6WB8WRJTCYCR1HB5M60VGZ8` · actor [claude-code-84def52f](../actor/claude-code-84def52f.md)

## result

PASS

## statement

QA cycle 2 at 76c40f51, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows plan/bugfix/other in .rhiza/workflows, RHIZA_* unset): run ls --json shows pick 'rule request.bug' for the bugfix plan run and 'explicit' for --workflow plan runs; run start --dry-run prints 'Picked by:' with rule ready.bug, request or explicit; the daemon's start lines name the pick ('picked by rule ready.feature', 'picked by default').

## Links

### Out

- verifies [run-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/run-records-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
