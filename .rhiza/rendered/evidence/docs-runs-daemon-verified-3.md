# docs-runs-daemon-verified-3

evidence · `01M6ZKAT5RYF2RH6E4Z3YGHSR0` · actor [claude-code-b108cb8e](../actor/claude-code-b108cb8e.md)

## result

PASS

## statement

Verified docs/runs-and-daemon.md at d7d339c1 as verifier claude-code-b108cb8e (attempt 5). Built this head into my scratchpad and, with every RHIZA_* var unset and a hermetic HOME, followed the page: `rhiza root init <tmp>/code`, `cd <tmp>/code`, `rhiza project add <tmp>/code/myapp` now succeeds (attempt 4's finding fixed). `rhiza run start --workflow ship-change --dry-run` without --permission-mode stops asking for one, as the page says; with --permission-mode acceptEdits it passes that check (stops only on a leafless test change set); bypassPermissions is refused naming trust = "dedicated" in machine.toml, matching internal/api/session.go bypassCheck. `rhiza workflow new mine --from ship-change`, `workflow check mine` and `run start --workflow mine --permission-mode acceptEdits --dry-run` all exit 0. Adding an [autonomy] table with signoff/ship = owner makes the dry-run print 'Autonomy ship: owner from workflow mine' while plan/accept come from rhiza.autonomy.*, matching the Per workflow section; `rhiza daemon settings` shows defaults plan=owner, accept/ship/signoff=auto, matching the Autonomy section. The rhiza.daemon.permissionMode wording matches daemon.go (queued starts fall back to the daemon mode, but run start requires its own). request-changes now shows <findings.json> <run>. docs-links-resolve: `go test -trimpath -count=1 ./internal/doclinks/...` passes. docs-commands-exist has no check built yet, so by hand every `rhiza <cmd> [<sub>]` named on the page answers --help, and the flags it names (--workflow, --permission-mode, --foreground, --autonomous, --kind, --ref, --note) exist. Every earlier confirmed finding (autonomy defaults, per-workflow autonomy, root setup, request-changes, project add cwd, permission mode) is fixed. Unconfirmed nits left: 'one per repository' (the daemon is per Rhiza root) and 'Set them yourself if you want to be asked' above two auto lines. Prose reads plain, consistent with doc-prose-humanized; the owner read is left to the owner at review.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
