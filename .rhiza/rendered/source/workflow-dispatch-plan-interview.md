# workflow-dispatch-plan-interview

source · `01M6VYSB8HVGWDQP0DECJFHSY9`

## quotes

```json
[]
```

## raw_sha256

35581b860ad363166917cfd095b793cc2cd193e4fc59e313d0da906a53bdc665

## session

b640ca70-18c2-43e6-a7ec-01a1b7e43901

## summary

Autonomous plan run on change set 01M6VY7G (Workflow dispatch). No owner took part: the planner raised two batches of questions and the arbiter answered each with a Decision worded (assumed, would have escalated). Facts from the code: the daemon starts every request's plan run with the built-in plan workflow (internal/api/request.go daemonStartPlan, planWorkflow in run.go), and every ready change set gets ship-change (defaultWorkflow in run.go); run start --new refuses any workflow other than plan; the change set records no request kind today; rhiza.toml is the root's machine-wide config (internal/rootcfg), shared by several projects and not committed in any repo. Settled by the arbiter in workflow-dispatch-shape: the change set serves milestone-25, the backlog milestone; the rules live in a committed repo file, .rhiza/dispatch.toml, not in rhiza.toml; rhiza request --workflow overrides only the ship pick and is recorded on the change set; a manual run start with no --workflow resolves the same way as the daemon. Settled in workflow-dispatch-rules-file: one table per trigger ([request] and [ready]) keyed by kind, with missing keys falling back to the built-in defaults and any other table or key an error; an unresolvable workflow name is a broken reference (rhiza check reports it, request --workflow refuses it, the daemon starts nothing and raises attention, never a silent fallback); a change set no request opened has kind feature; each run records what dispatch picked and why, shown in run ls and run status. Also found: rhiza change show fails on this change set because model.ShowFieldValue re-encodes a field with an @-reference using json.Marshal, which escapes '>' and fails the canonical check. That is outside this intent and goes to its own draft change set. Decision workflows-are-the-rules (PR #353) was not landed at plan time, so it is cited here but not linked.

## title

Workflow dispatch plan, autonomous run with the arbiter

## Links

### In

- derived-from [workflow-dispatch-shape](../decision/workflow-dispatch-shape.md)
- derived-from [workflow-dispatch-rules-file](../decision/workflow-dispatch-rules-file.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
