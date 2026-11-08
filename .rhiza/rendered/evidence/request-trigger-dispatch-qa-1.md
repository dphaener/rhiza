# request-trigger-dispatch-qa-1

evidence · `01M6W8JMT8HFGQRQ7A20AJ0ZDE` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): With [request] bug = "other", daemon --once logged 'starting the plan run of request ... with workflow other (picked by rule request.bug)' and detached 'run start --foreground --workflow other --pick rule request.bug --new ... --request ...', which ran to succeeded. Feature requests in the same polls started with plan (picked by default). With [request] milestone = "other" and daemon --autonomous --permission-mode bypass, the milestone request started as '--workflow other ... --milestone' with no --autonomous, and its run opened a record change set.

## Links

### Out

- verifies [request-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/request-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
