# ready-trigger-dispatch-qa-2

evidence · `01M6WB8WBD3EEC0P2NAZJ6TP55` · actor [claude-code-84def52f](../actor/claude-code-84def52f.md)

## result

PASS

## statement

QA cycle 2 at 76c40f51, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows plan/bugfix/other in .rhiza/workflows, RHIZA_* unset): run start --dry-run on bug change sets: with no file the default ship-change ('has no leaves: ship-change builds...'); with [ready] bug = "other" 'Workflow: other, Picked by: rule ready.bug'; --workflow bugfix gives 'Picked by: explicit'; a change set whose request named bugfix gives 'Workflow: bugfix, Picked by: request' even beside a broken file. With '[ready] bug = 3' the start now exits 1 with the file's load error and never falls back to ship-change; --workflow other still starts.

## Links

### Out

- verifies [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
