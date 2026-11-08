# review-01m6vy7e-r3-f3

question · `01M6VZP9JVHFVPNR4SRSMN2NMT`

## statement

Plan review round 3 left finding f3 (contradiction on workflow-dispatch-rules-file): The decision (and leaf dispatch-rules-file) lets [ready] map the kind milestone to a workflow, but the same plan makes that key unreachable: a milestone request's plan run opens a record change set, which lands as is and is never accepted and ready (changeset fold.go:606; request-trigger-dispatch 'a milestone request still opens a record change set'), and change sets a milestone plan opens get kind feature (decision: 'A change set no request opened ... has kind feature'). So `[ready] milestone = "x"` passes rhiza check yet can never fire. Either drop milestone from [ready]'s keys (an error, like any other unknown key) or say what change set ever has kind milestone at the ready trigger. The critic confirmed it: The ready trigger fires on accepted, ready change sets that serve a milestone (internal/daemon/pick.go:13,84). A milestone request's plan run opens a record change set (request.go:43, run.go:703) which lands straight from in-flight (fold.go:606, lifecycle.go:160) and is never accepted-and-ready in the normal flow, and the decision gives every change set no request opened kind feature. So no change set reaches the ready trigger with kind milestone, yet the decision and dispatch-rules-file accept [ready] milestone as a valid key that rhiza check passes. The plan must drop the key or say which change set carries kind milestone at ready. How should it be settled?

## Links

### In

- answers [dispatch-ready-no-milestone-key](../decision/dispatch-ready-no-milestone-key.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
