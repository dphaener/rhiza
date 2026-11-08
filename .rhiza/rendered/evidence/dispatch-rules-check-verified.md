# dispatch-rules-check-verified

evidence · `01M6W356SBNK1R4ZQ8E96YMM77` · actor [claude-code-77e5de6d](../actor/claude-code-77e5de6d.md)

## result

PASS

## statement

At 907c325a, go test -trimpath ./cmd/rhiza -run 'TestScripts/dispatch-check$' -v passes: the script runs rhiza check passing with no .rhiza/dispatch.toml, with [request] bug naming a repo workflow in .rhiza/workflows/bugfix.toml, and with [ready] feature naming the built-in ship-change; and failing with stdout 'dispatch.toml: request.bug names workflow "nothing-here"', with 'dispatch.toml: ready.milestone names workflow "personal-only"' while that workflow exists only in $XDG_CONFIG_HOME/rhiza/workflows set inside the script, and with 'unknown table [reqest]'. describe-golden passes with the new dispatch rule; internal/check, internal/dispatch and the api Check/Dispatch/Workflow tests pass. Diff read: dispatchProblems in internal/api/check.go resolves names only by os.Stat under .rhiza/workflows and workflow.Builtin, never the personal layer; findings carry rule dispatch, path .rhiza/dispatch.toml and severity fail. The full internal/api package went past a 240s timeout under load with no test failing, a pre-existing slowness and not this diff.

## Links

### Out

- verifies [dispatch-rules-check](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-check.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
