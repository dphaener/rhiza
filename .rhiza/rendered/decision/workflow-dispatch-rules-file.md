# workflow-dispatch-rules-file

decision · `01M6VYQWKXT7DC11ET2AAAER7G`

## reasoning

Builds on @workflow-dispatch-shape (rules in a committed .rhiza/ file, built-in defaults, override targets the ship pick). Shape: the intent has exactly two triggers and the request kinds already fixed in internal/api/request.go (enum feature,bug,milestone), so a table per trigger keyed by kind is the smallest file that covers every case; partial files falling back key by key keep 'built-in defaults reproduce today's behaviour'. Bad names: decision workflows-are-the-rules says Rhiza enforces only references between nodes, and a rule naming a missing workflow is exactly a broken reference, so `rhiza check` (run by make check and CI) catches it before it lands; a silent fallback to ship-change would run the wrong workflow on the owner's work, so runtime failures surface as attention instead. Validating --workflow at queue time reuses resolveWorkflow (internal/api/run.go) and fails fast while the owner is still there. Default kind: change sets from milestone plans and `rhiza change new` are new work, which is what feature means and what the default --kind already is (request.go: 'default feature'). Visibility: run journals already carry the workflow name (run.go Workflow field, run_ls.go), so adding the reason there puts it where the owner's Rex apps and `rhiza run ls` already read, without a new surface.

## statement

.rhiza/dispatch.toml has one table per trigger, [request] (a request is queued) and [ready] (a change set is accepted and ready), and each maps a kind (feature, bug, milestone) to a workflow name, e.g. `[ready] bug = "bugfix"`. A key left out falls back to the built-in default (request -> plan, ready -> ship-change); any other table or key is an error. A workflow name no layer resolves is a broken reference: `rhiza check` reports a rule in .rhiza/dispatch.toml whose workflow neither the repo (.rhiza/workflows/) nor the built-in layer resolves (it never reads the personal layer, which differs per machine and is absent in CI, so the same committed .rhiza/ checks the same everywhere), `rhiza request --workflow` refuses it with a usage error when queued, and if a name stops resolving by the time the daemon dispatches, the daemon starts no run and raises it as attention; it never silently falls back to the default. A change set no request opened (from a milestone plan or `rhiza change new`) has kind feature; one a request's plan run opens records that request's kind. Each run records which workflow dispatch picked and why (the request's override, the rule trigger.kind, or the built-in default), shown in `rhiza run ls` and `rhiza run status` and logged by the daemon; `rhiza change show` shows the change set's recorded kind and workflow.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b640ca70-2](../question/question-b640ca70-2.md)
- derived-from [workflow-dispatch-plan-interview](../source/workflow-dispatch-plan-interview.md)

### In

- decided-by [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)
- decided-by [dispatch-rules-check](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-check.md)
- decided-by [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)
- decided-by [request-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/request-workflow-flag.md)
- decided-by [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)
- decided-by [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)
- decided-by [dispatch-unresolved-attention](../capability/rhiza/workflows/workflow-dispatch/dispatch-unresolved-attention.md)
- decided-by [run-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/run-records-dispatch.md)
- decided-by [dispatch-guide-topic](../capability/rhiza/workflows/workflow-dispatch/dispatch-guide-topic.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
