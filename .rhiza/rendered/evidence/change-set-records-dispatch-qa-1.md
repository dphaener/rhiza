# change-set-records-dispatch-qa-1

evidence · `01M6W8JKVY5G80KVNTEABM3CW1` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): rhiza change new gives change show --json kind intention, request_kind feature, workflow null, and text '- Request kind: feature' with no Workflow line. The bug request queued with --workflow bugfix, claimed by daemon --once and run by the logged run start --foreground --new ... --request, opened 'Fix Crash on empty input' with request_kind bug and workflow bugfix (text: 'Request kind: bug', 'Workflow: bugfix'). A milestone request's run opened a record change set with request_kind milestone and no workflow. rhiza check stays ok.

## Links

### Out

- verifies [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
