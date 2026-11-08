# watchdog-green-pr-unmerged-verified

evidence · `01M6WX5395J6G27FPJX78CRWF6` · actor [claude-code-e5c4dacd](../actor/claude-code-e5c4dacd.md)

## result

PASS

## statement

At 93c41166, go test -trimpath -count=1 ./internal/api -run TestWatchdogGreenUnmerged -v passes TestWatchdogGreenUnmerged and TestWatchdogGreenUnmergedReadsOnce. With a fake gh read (watchdogDeps.ci) and no merge dep, the pass resumes (--retry via queueContinue) only the autonomous run failed at its merge step, leaves the live run at merge alone, and holds the succeeded-unmerged run and the run gated by a signoff as inbox items with no resume and no merge call; red/absent checks or no open PR do nothing, and a spent budget holds. Attempt 1's finding (gh read per ended run every poll) is fixed: a run found with no open PR is remembered as PRGone at its last journal entry, persisted in daemon state, and skipped until the journal moves; the ReadsOnce test shows 1 read over three polls through watchdogPasses and one more after the journal moved. All TestWatchdog* in internal/api and internal/daemon tests pass; go vet and gofmt clean.

## Links

### Out

- verifies [watchdog-green-pr-unmerged](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-green-pr-unmerged.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
