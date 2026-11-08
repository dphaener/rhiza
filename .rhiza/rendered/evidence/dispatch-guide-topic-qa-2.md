# dispatch-guide-topic-qa-2

evidence · `01M6WB9752XWNFXGM2T3C1BBP3` · actor [claude-code-84def52f](../actor/claude-code-84def52f.md)

## result

PASS

## statement

QA cycle 2 at 76c40f51, built binary: rhiza guide workflows prints the Dispatch section naming .rhiza/dispatch.toml, [request] (feature, bug, milestone) and [ready] (now feature, bug only), the defaults plan and ship-change, rhiza request --workflow, the ready order (run start --workflow, the request's workflow, the [ready] rule, ship-change), and the workflow-unresolved skip, which matches what daemon status and inbox showed hands-on. go test -trimpath ./internal/guide/... passes. The stray 'opens,.' is still there.

## Links

### Out

- verifies [dispatch-guide-topic](../capability/rhiza/workflows/workflow-dispatch/dispatch-guide-topic.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
