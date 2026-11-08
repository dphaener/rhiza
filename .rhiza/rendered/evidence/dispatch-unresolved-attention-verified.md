# dispatch-unresolved-attention-verified

evidence · `01M6W6W9W9JHDVDFKF9G5VJVJK` · actor [claude-code-817af69f](../actor/claude-code-817af69f.md)

## result

PASS

## statement

Verifier claude-code-817af69f at 8bdad8b6: read the diff (internal/api daemon.go, request.go, inbox.go, internal/daemon/pick.go) and ran go test -trimpath -count=1 ./internal/api -run 'Unresolved|Dispatch|PlanRequests|DaemonStatus|Inbox' (ok, run from a scratch worktree outside ~/code; inside ~/code the rootguard flags the live root's state folder, a known environmental issue) plus go test ./cmd/rhiza TestScripts describe-golden/daemon*/inbox*/request*/dispatch* (ok). TestDaemonSkipsReadyWithUnresolvedWorkflow: a daemon poll with a [ready] rule naming the missing workflow 'ghost' picks nothing (no default fallback), records a workflow-unresolved skip with Detail [ghost, the [ready] rule for feature in .rhiza/dispatch.toml, error], workflowItems lists one 'workflow' inbox item naming ghost and the rule and no ID; after the rule is removed the next poll starts it. TestDaemonSkipsRequestWithUnresolvedWorkflow: a queued request whose recorded workflow is 'ghost' starts no plan run and gets a skip (request ID, text, ghost, 'named by the request (--workflow)') and an inbox item, while the other request in the same daemonPlanRequests pass still starts. The poll merges request skips into daemon state, so daemon status and the inbox and notices show them.

## Links

### Out

- verifies [dispatch-unresolved-attention](../capability/rhiza/workflows/workflow-dispatch/dispatch-unresolved-attention.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
