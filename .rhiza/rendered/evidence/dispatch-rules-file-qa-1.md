# dispatch-rules-file-qa-1

evidence · `01M6W8JJBM56DT3GDVV9FE4J4P` · actor [claude-code-ed5eb908](../actor/claude-code-ed5eb908.md)

## result

FAIL

## statement

QA at 75f0a821, built binary in a scratch Rhiza root (fake rex on PATH, one-step command workflows, RHIZA_* unset): rhiza check with .rhiza/dispatch.toml: unknown table [reqest] fails 'unknown table [reqest]: want [request] or [ready]'; [request] chore fails 'unknown key "chore"'; bug = 3 fails 'want a workflow name as a string, got int64'; bug = "../x" fails 'not a bare workflow name'; a broken TOML fails naming the line. With no file, run start --dry-run picks ship-change (default) and the daemon starts requests with plan (picked by default). But [ready] milestone = "bugfix" passes rhiza check (ok), against decision dispatch-ready-no-milestone-key, which makes it an unknown-key error.

## Links

### Out

- verifies [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
