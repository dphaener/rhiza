# request-run-landed-resume-check

decision · `01M6VZY5Q97N5MKT6RKN0W7181`

## reasoning

Read internal/api/run.go: run resume exempts a landed change set only when st.Workflow == planWorkflow or the run already visited land; otherwise Terminal() returns a Conflict. The leaf's own statement says every check that asks 'is this a request's plan run' moves off the name plan, and the early-landed record is exactly why plan is exempt here, so this check is the same test and was simply omitted. Leaving it out makes daemon.Resumable pick a run that run resume then refuses, which defeats the leaf's intent (dispatching request kinds to other workflows). The critic confirmed the gap. Extending the list and acceptance stays inside the change set's intent and adds no new mechanism; asserting the resume itself, not just Resumable's pick, is what would have caught it.

## statement

@request-trigger-dispatch also rekeys run resume's landed-state check (internal/api/run.go, the `st.Workflow == planWorkflow && StateLanded` exemption beside the visited-land one) on 'is a request's plan run' (started with --new or for a request) rather than the workflow name plan, so a request run under any dispatched workflow whose change set landed early is let through rather than refused as terminal. The leaf's acceptance tests resume such a run end to end through run resume under a workflow not named plan, not only see it in daemon.Resumable's list

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6vy7e-r3-f2](../question/review-01m6vy7e-r3-f2.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
