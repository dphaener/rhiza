# review-01m6vy7e-dispatch-r3-f1

question · `01M6VZP9EAZX6GB0QF2R2WHCXP`

## statement

Plan review round 3 left finding f1 (contradiction on change-set-records-dispatch): The leaf has a change set 'record its kind' (feature, bug, milestone) and its acceptance checks 'kind bug' and 'kind feature' in change show --json. But a change set already has a kind with another meaning: changeset.Kind is intention or record (internal/changeset/body.go:66-77), `rhiza change new --kind` takes intention|record (internal/api/change.go:63), and change show --json already prints "kind": "intention"/"record". As written, the new field collides with or overwrites the existing one, and a milestone request's change set would need kind both 'record' and 'milestone'. The leaf (and decision workflow-dispatch-rules-file, 'rhiza change show shows the change set's recorded kind') should name a distinct field, e.g. request_kind, in the open entry and in change show. The critic confirmed it: A change set already has a kind with another meaning: changeset.Kind is intention|record (internal/changeset/body.go:66-77), change new --kind takes intention|record (internal/api/change.go:63) and change show --json prints "kind" from c.Open.Kind (change.go:220,304). change-set-records-dispatch says the change set 'records its kind' and its acceptance reads 'kind bug'/'kind feature' in change show --json, so the new field collides with the existing one (a milestone request's change set is kind record). The leaf and decision workflow-dispatch-rules-file need a distinct field name. How should it be settled?

## Links

### Out

- blocks [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)

### In

- answers [dispatch-kind-field-request-kind](../decision/dispatch-kind-field-request-kind.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
