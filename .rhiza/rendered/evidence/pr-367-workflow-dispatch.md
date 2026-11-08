# pr-367-workflow-dispatch

evidence · `01M6WFZEGK506MM9NMK3G5W2PV` · actor [dphaener](../actor/dphaener.md)

## head

42811a775932dab046118fa54b5793d17047d74e

## pr

https://github.com/dphaener/rhiza/pull/367

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #367 (https://github.com/dphaener/rhiza/pull/367) passed CI at 42811a7

## status

verified

## Links

### Out

- verifies [workflow-dispatch](../capability/rhiza/workflows/workflow-dispatch.md)
- verifies [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)
- verifies [dispatch-rules-check](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-check.md)
- verifies [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)
- verifies [request-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/request-workflow-flag.md)
- verifies [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)
- verifies [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)
- verifies [dispatch-unresolved-attention](../capability/rhiza/workflows/workflow-dispatch/dispatch-unresolved-attention.md)
- verifies [run-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/run-records-dispatch.md)
- verifies [dispatch-guide-topic](../capability/rhiza/workflows/workflow-dispatch/dispatch-guide-topic.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
