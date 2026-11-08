# ready-trigger-dispatch-qa-1

evidence · `01M6W8JN8PKB03BTE7DXZ4Q4EK` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

PASS

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): run start --dry-run --json: the bug change set whose request recorded workflow bugfix picks bugfix (pick request) even with [ready] bug = "other"; --workflow plan wins (pick explicit); a feature change set with [ready] feature = "other" picks other (rule ready.feature); with no dispatch file it picks ship-change (default), and the text output shows 'Workflow: ship-change' / 'Picked by: default'. After the recorded bugfix workflow file was removed, run start on that change set exits 2 'no workflow "bugfix"' and does not fall back.

## Links

### Out

- verifies [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
