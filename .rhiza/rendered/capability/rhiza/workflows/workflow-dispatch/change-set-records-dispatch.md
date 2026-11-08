# change-set-records-dispatch

capability · `01M6VYTAJBDW1YH8D7Y1XGCK76` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript shows a change set opened by run start --new --request for a queued bug request whose queue entry carries a workflow (the entry written directly into the request queue file by the script, since request --workflow arrives later in request-workflow-flag) recording kind bug and that workflow in change show --json, a change set from rhiza change new showing kind feature and no workflow, and the log still folding with fold-ok

## statement

A change set records its kind and an optional workflow: the change set a request's plan run opens records that request's kind and the workflow the request named, any other change set (from a milestone plan or rhiza change new) has kind feature and no workflow, and rhiza change show shows both

## Links

### Out

- decided-by [workflow-dispatch-rules-file](../../../../decision/workflow-dispatch-rules-file.md)
- decided-by [workflow-dispatch-shape](../../../../decision/workflow-dispatch-shape.md)

### In

- depends-on [request-workflow-flag](request-workflow-flag.md)
- depends-on [ready-trigger-dispatch](ready-trigger-dispatch.md)
- blocks [review-01m6vy7e-dispatch-r3-f1](../../../../question/review-01m6vy7e-dispatch-r3-f1.md)
- verifies [change-set-records-dispatch-verified](../../../../evidence/change-set-records-dispatch-verified.md)
- verifies [change-set-records-dispatch-qa-1](../../../../evidence/change-set-records-dispatch-qa-1.md)
- verifies [change-set-records-dispatch-qa-2](../../../../evidence/change-set-records-dispatch-qa-2.md)
- verifies [pr-367-workflow-dispatch](../../../../evidence/pr-367-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
