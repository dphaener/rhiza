# dispatch-rules-check-verified-2

evidence · `01M6W9KXVX02JYWA4KA2PH2M8Q` · actor [claude-code-23e9c857](../actor/claude-code-23e9c857.md)

## result

PASS

## statement

Verified at a97778b2 (attempt 2). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/dispatch-check$' -count=1: ok. The script, with XDG_CONFIG_HOME set to $WORK/xdg inside it, shows rhiza check passing with no dispatch.toml, with [request] bug naming .rhiza/workflows/bugfix.toml, and with [ready] feature naming built-in ship-change. It fails with request.bug naming missing "nothing-here" (stdout names request.bug), with ready.feature naming "personal-only", which exists only in $XDG_CONFIG_HOME/rhiza/workflows (stdout names ready.feature), with [ready] milestone (unknown key), and with unknown table [reqest] (stdout names the table). Compliance f2 is fixed: the personal-layer case now uses the valid [ready] feature key, and dispatchProblems (internal/api/check.go:530) iterates dispatch.TableKinds per table and resolves only the repo layer and built-ins, never the personal layer. go test -trimpath ./internal/dispatch/... also passes: ok.

## Links

### Out

- verifies [dispatch-rules-check](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-check.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
