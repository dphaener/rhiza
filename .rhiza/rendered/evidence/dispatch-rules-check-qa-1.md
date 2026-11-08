# dispatch-rules-check-qa-1

evidence · `01M6W8JK18DNGZY33A72FMKA73` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): rhiza check passes with no dispatch file, with [request] bug naming the repo workflow .rhiza/workflows/bugfix.toml and with [ready] bug = ship-change (built in). It fails (exit 1) for [request] bug = "nothing" with 'fail dispatch .rhiza/dispatch.toml: ... request.bug names workflow "nothing", which is neither in .rhiza/workflows nor built in', and for [ready] bug = "pers" where pers exists only in $XDG_CONFIG_HOME/rhiza/workflows ('ready.bug names workflow "pers"'), so the personal layer is not read.

## Links

### Out

- verifies [dispatch-rules-check](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-check.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
