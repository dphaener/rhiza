# dispatch-ready-no-milestone-key

decision · `01M6W005N7EKNHZQB6D7XNGEZB`

## reasoning

No change set ever reaches the ready trigger with request_kind milestone: a milestone request's plan run opens a record change set, which lands straight from in-flight and is never accepted-and-ready (changeset fold.go:606, lifecycle.go:160; request-trigger-dispatch keeps that), and change sets a milestone plan or rhiza change new opens get request_kind feature (change-set-records-dispatch, dispatch-kind-field-request-kind). A [ready] milestone rule would pass rhiza check yet never fire, which is the silent no-op the plan otherwise refuses (unknown keys are errors, unresolved workflows are attention). Dropping the key is the smaller change and keeps the file honest; inventing a path that gives a ready change set kind milestone would add lifecycle behaviour the intent never asked for. If such a path is added later, the key can be admitted then.

## statement

The [ready] table in .rhiza/dispatch.toml accepts only the kinds feature and bug; a milestone key under [ready] is an error naming it, like any other unknown key, and rhiza check reports it. [request] keeps all three kinds (feature, bug, milestone). This narrows workflow-dispatch-rules-file and dispatch-rules-file: 'each maps a kind (feature, bug, milestone)' holds for [request] only, and the ready resolver never looks up a milestone rule.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6vy7e-r3-f3](../question/review-01m6vy7e-r3-f3.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
