# dispatch-unresolved-attention

capability · `01M6VYTAYZW0RTPZ31DVTTYVTB` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

Tests in internal/api show a daemon poll with a recorded workflow or a rule naming a workflow that does not resolve starting no run, recording a skip naming the workflow and its source, and listing an inbox item for it, while other candidates in the same poll still start

## statement

When the workflow dispatch picks for a request or a ready change set no longer resolves at start time, the daemon starts no run, never falls back to the default, and raises it as attention: a skip reason in daemon status and an inbox item naming the change set or request, the workflow and where the pick came from

## Links

### Out

- depends-on [request-trigger-dispatch](request-trigger-dispatch.md)
- depends-on [ready-trigger-dispatch](ready-trigger-dispatch.md)
- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)

### In

- verifies [dispatch-unresolved-attention-verified](../../../../evidence/dispatch-unresolved-attention-verified.md)
- verifies [dispatch-unresolved-attention-qa-1](../../../../evidence/dispatch-unresolved-attention-qa-1.md)
- verifies [dispatch-unresolved-attention-qa-2](../../../../evidence/dispatch-unresolved-attention-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
