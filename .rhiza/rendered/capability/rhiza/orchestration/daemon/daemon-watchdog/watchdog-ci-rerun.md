# watchdog-ci-rerun

capability · `01M6WJ2CXX5VYQ2236KG5ZEX7X` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogCIRerun, with a fake gh, shows the failed jobs rerun once, the run resumed after they pass, no second rerun on the same head, no rerun of a head whose run attempt is already above 1, and a new head allowed its own rerun

## statement

When a run ended failed waiting on CI with red checks, the daemon reruns the PR's failed jobs once per PR head and resumes the run with --retry when they pass. A head whose failed jobs were already rerun (gh shows a run attempt above 1, as the plan step's own rerun leaves it) counts as that one rerun, so the watchdog does not rerun it again; a second red on the same head is left for the owner

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-ci-rerun-verified](../../../../../evidence/watchdog-ci-rerun-verified.md)
- verifies [watchdog-ci-rerun-verified-3](../../../../../evidence/watchdog-ci-rerun-verified-3.md)
- verifies [watchdog-ci-rerun-qa-2](../../../../../evidence/watchdog-ci-rerun-qa-2.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
