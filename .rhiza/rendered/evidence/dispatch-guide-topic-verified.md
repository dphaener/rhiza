# dispatch-guide-topic-verified

evidence · `01M6W6KNWY80T8SV2KPD8MD6A3` · actor [claude-code-2a4bf194](../actor/claude-code-2a4bf194.md)

## result

PASS

## statement

At 3866a47a ran go test -trimpath -count=1 -v ./internal/guide/...: all 20 tests PASS, including TestWorkflowsText, whose new dispatch group requires .rhiza/dispatch.toml, [request], [ready], rhiza request --workflow, feature/bug/milestone, ship-change and workflow-unresolved in the workflows topic, and the existing size checks TestCaps, TestCapCatchesOversize and TestTokens. Read the diff: the new Dispatch section of internal/guide/text/workflows.md matches internal/dispatch (two tables, three kinds, an error naming any other table, key or value, defaults plan and ship-change, no kind means feature) and readyPick/startWorkflow in internal/api/run.go (explicit --workflow, then the request's workflow, then the [ready] rule, then ship-change). It also covers the workflow-unresolved skip with no fallback. Minor: a stray comma in 'the request opens,.' for the review run to fix.

## Links

### Out

- verifies [dispatch-guide-topic](../capability/rhiza/workflows/workflow-dispatch/dispatch-guide-topic.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
