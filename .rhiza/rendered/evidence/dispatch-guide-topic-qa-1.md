# dispatch-guide-topic-qa-1

evidence · `01M6W8JQ64SN394R0ZREJ0C29C` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): rhiza guide workflows prints a Dispatch section naming .rhiza/dispatch.toml, [request], [ready], feature/bug/milestone, the defaults plan and ship-change, rhiza request --workflow, the ready order (run start --workflow, the request's workflow, the [ready] rule, ship-change) and the workflow-unresolved skip, which matches daemon.SkipWorkflow. The stray comma 'the request opens,.' (line 63) is still there.

## Links

### Out

- verifies [dispatch-guide-topic](../capability/rhiza/workflows/workflow-dispatch/dispatch-guide-topic.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
