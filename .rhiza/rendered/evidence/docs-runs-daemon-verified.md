# docs-runs-daemon-verified

evidence · `01M6ZHJ6ZQZGQEFQ4F08GDHBR5` · actor [claude-code-ea0f5f18](../actor/claude-code-ea0f5f18.md)

## result

PASS

## statement

At eb28791e read docs/runs-and-daemon.md (184 lines): it covers Rex and Claude Code as requirements (rex install, guide install), workflows and runs, sign-off, questions/inbox, the daemon and its settings, requests and .rhiza/dispatch.toml [request]/[ready], per-stage autonomy (plan/accept/signoff/ship), [[schedules]] (every/cron/CRON_TZ/workflow/command/catch_up/overlap) and the bypassPermissions risk warnings. go test -trimpath ./internal/doclinks/... passes (TestRepoDocsLinksResolve over the real docs, incl. agents.md and concepts/change-sets.md). docs-commands-exist has no check in the tree yet, so with a scratch build I ran --help on all 28 rhiza commands the doc names (rex install, workflow ls/new/describe/check, run start/ls/status/tail/resume/abandon/approve/request-changes/send-back, inbox, question answer, daemon start/status/stop/install/settings/set, health, request, request ls, guide install/workflows, describe): all exit 0 (an unknown command exits 2). Flags named exist (run start --workflow/--foreground/--autonomous, global --dry-run, request --kind/--ref, rex install --force, workflow new --from, run send-back --note). Settings keys match daemon settings output; dispatch defaults plan/ship-change, schedule keys, 1m minimum and CRON_TZ match internal/dispatch and internal/schedule; the warnings match internal/api/request_exposure.go (health and workflow check warn, nothing refuses). Prose reads in the owner's plain voice. Owner read-through is left to the owner.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
