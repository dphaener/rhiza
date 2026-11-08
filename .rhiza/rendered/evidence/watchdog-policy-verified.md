# watchdog-policy-verified

evidence · `01M6WQVCESZFKN6P09W7CD4TSX` · actor [claude-code-f2269798](../actor/claude-code-f2269798.md)

## result

PASS

## statement

At f73688bd I ran go test -trimpath -count=1 ./internal/api -run TestWatchdogPolicy -v: PASS, with subtests for the defaults (enabled, 20m/40m nudge/retry, 30m command, 90m wait, retryBudget 2, ciReruns 1, pollBehind 3, fileBugs and autonomousOnly true), a machine value under [daemon] in rhiza.toml, a rhiza.project.watchdog.* git config key over the machine key, a workflow step's stall field winning for that step, and invalid values rejected with their key named. go test -trimpath ./cmd/rhiza -run 'TestScripts/(watchdog-policy|describe-golden|daemon)': PASS. go test ./internal/workflow: PASS. With a scratch build in a fresh repo, rhiza daemon settings listed all 10 rhiza.daemon.watchdog.* keys as (default, machine) and the 9 rhiza.project.watchdog.* keys (no pollBehind) as (default, project). The settings, daemon-set and autonomy tests in internal/api pass. The whole internal/api package hit its 600s timeout at load 16, in an unrelated database/sql goroutine dump. I read the diff: the keys join settingDefs in init, the resolver takes project over machine over default, and stall is a positive step duration that a workflow-kind step rejects.

## Links

### Out

- verifies [watchdog-policy](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-policy.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
