# watchdog-ci-rerun-verified

evidence · `01M6WT95B0ZPVKG3DQ9MVB80DH` · actor [claude-code-a80772cc](../actor/claude-code-a80772cc.md)

## result

PASS

## statement

At f5b15eab, go test -trimpath ./internal/api -run TestWatchdogCIRerun -count=1 -v passed. The test fakes the gh-backed deps (ci reads head, checks and Actions run attempts; rerun records gh run rerun --failed calls) and shows: a red head's failed run 7 rerun once; still red at attempt 2 gives no second rerun and no resume; once green the run is resumed with --retry exactly once (two polls, one resume, one log line per action); a fresh red head whose run is already at attempt 2 is not rerun; a new head bbbb222 gets its own rerun. Read the diff: watchdogCIRerun only acts on ended, failed, non-live, autonomous runs whose failed step runs or waits on rhiza ship check (not --pr-open/--main), records CIHead before rerunning, and charges the resume to the watchdog budget; the live path reads gh pr view and gh run list --commit <head> --json databaseId,attempt. All TestWatchdog* tests and go vet on internal/api and internal/daemon pass.

## Links

### Out

- verifies [watchdog-ci-rerun](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-ci-rerun.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
