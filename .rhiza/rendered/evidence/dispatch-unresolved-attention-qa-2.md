# dispatch-unresolved-attention-qa-2

evidence · `01M6WB8WE4Q449B05HDCJ9DN8P` · actor [claude-code-84def52f](../actor/claude-code-84def52f.md)

## result

PASS

## statement

QA cycle 2 at 76c40f51, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows plan/bugfix/other in .rhiza/workflows, RHIZA_* unset): daemon --once on an accepted milestone change set: with '[ready] feature = 3' no run starts and daemon status shows workflow-unresolved with ready.feature, 'the [ready] rule for feature in .rhiza/dispatch.toml' and the load error; rhiza inbox lists 'change set "Built": no run started: ... does not resolve'. With feature = "nope" the skip names nope and its rule. Once feature = "other", the poll starts the run with --workflow other --pick rule ready.feature.

## Links

### Out

- verifies [dispatch-unresolved-attention](../capability/rhiza/workflows/workflow-dispatch/dispatch-unresolved-attention.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
