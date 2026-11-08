# review-01m6vy7e-r3-f2

question · `01M6VZP9FPP58AM3F8KM40FS24`

## statement

Plan review round 3 left finding f2 (contradiction on request-trigger-dispatch): The leaf lists the checks that must key on 'is a request's plan run' instead of the name plan: daemon.Resumable, the inbox, and run resume's milestone owner gate (run.go:849). It misses run resume's own landed-state check (internal/api/run.go:879-889): it lets a run on a landed change set through only if st.Workflow == planWorkflow or the run already visited land, else Terminal() refuses it with a Conflict. So a bug request dispatched to e.g. 'bugplan' whose change set lands early would be picked by daemon.Resumable (as the acceptance tests) and then refused by run resume, stuck. The statement and acceptance should include run.go:879 (and the acceptance should resume the run, not just see it in Resumable's list). The critic confirmed it: run resume (internal/api/run.go:879-889) lets a non-ended run on a landed change set through only if st.Workflow == planWorkflow or the run already visited land; otherwise Terminal() refuses it with a Conflict. request-trigger-dispatch lists daemon.Resumable, the inbox and the milestone owner gate (run.go:849) as the checks to rekey, but not this one, and its acceptance only checks Resumable's pick, so a request run under a non-plan workflow whose change set landed early would be picked and then refused by run resume. How should it be settled?

## Links

### Out

- blocks [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)

### In

- answers [request-run-landed-resume-check](../decision/request-run-landed-resume-check.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
