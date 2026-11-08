# dispatch-unresolved-attention-qa-1

evidence · `01M6W8JPBDNFC585F0AETE7SW5` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): With [request] bug = "ghost" (no such workflow), daemon --once started no run for the 'Ghost bug' request and still started the other queued feature requests in the same poll. daemon status shows 'Skipped: request ... Ghost bug: its workflow does not resolve: no run started (ghost; the [request] rule for bug in .rhiza/dispatch.toml; ...)', and rhiza inbox lists a workflow item 'request "Ghost bug": no run started: workflow ghost, the [request] rule for bug ... does not resolve'. The ready-change-set path was not reached by hand: the daemon skipped that change set first for serving no milestone.

## Links

### Out

- verifies [dispatch-unresolved-attention](../capability/rhiza/workflows/workflow-dispatch/dispatch-unresolved-attention.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
