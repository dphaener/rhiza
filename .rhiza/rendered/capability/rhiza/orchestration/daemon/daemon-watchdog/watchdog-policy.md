# watchdog-policy

capability · `01M6WJ2BR1XZB2E5F21H8TNBFH` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogPolicy covers the defaults with no watchdog key set, a machine value under [daemon] in rhiza.toml, a rhiza.project.watchdog.* key in git config giving that project's effective value over the machine key, a workflow step's stall field winning for that step, and an invalid value rejected with its key named; go test -trimpath ./cmd/rhiza -run TestScripts/watchdog-policy shows the machine and project keys with their sources and scopes in rhiza daemon settings

## statement

The daemon reads its watchdog policy from settings declared in settingDefs like the other daemon settings, each key of one scope: machine keys rhiza.daemon.watchdog.* as flat full-key strings under [daemon] in the root's rhiza.toml, and project keys rhiza.project.watchdog.* of the same names (all but pollBehind) in the project's git config. For a project the effective value is its project key when set, else the machine key, else the default. A workflow step may set its own stall limit as a step field, which wins for that step. Shipped defaults: enabled = true; agentNudge = 20m and agentRetry = 40m for agent steps; commandLimit = 30m; waitLimit = 90m for CI and PR wait steps; retryBudget = 2 per run; ciReruns = 1 per PR head; pollBehind = 3 (times the poll interval); fileBugs = true; autonomousOnly = true. rhiza daemon settings lists every machine and project key with its value, source and scope

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)
- decided-by [workflows-are-the-rules](../../../../../decision/workflows-are-the-rules.md)
- decided-by [warn-dont-safeguard](../../../../../decision/warn-dont-safeguard.md)

### In

- verifies [watchdog-policy-verified](../../../../../evidence/watchdog-policy-verified.md)
- verifies [watchdog-policy-qa-1](../../../../../evidence/watchdog-policy-qa-1.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
