# watchdog-ci-rerun-verified-3

evidence · `01M6X0MJRQBEARYS0D3D4T4F8P` · actor [claude-code-d272dcf4](../actor/claude-code-d272dcf4.md)

## result

PASS

## statement

Verifier claude-code-d272dcf4 at a1a714e0 (attempt 3, after review quality f1): go test -trimpath ./internal/api -run TestWatchdogCIRerun -v -count=1 PASS. Its fake gh shows the failed run 7 rerun once on a red head, nothing while the rerun runs, one resume with --retry after green and none on the next poll, no second rerun and no resume when the same head goes red again after the rerun (CIDone set), no rerun of a head whose run is already at attempt 2, and a rerun for a new head. TestWatchdogCIReadsBounded (via watchdogPasses) shows each run's PR read at most once per poll across the three PR passes: over three polls, gone-PR and red-after-rerun runs are read once, a run with its budget spent is never read, and a run is read again once its journal moves on. go test -run Watchdog ./internal/api, go test ./internal/daemon/... and go vet ./internal/api all pass. In the diff of 47741152, watchdogReadOnce wraps d.ci per poll; watchdogPRIdle skips runs while the journal end equals PRGone or CIDone, or when the budget is spent; rerunSpent replaces the old Attempts loop with the same meaning. Trade-off noted, as the f1 finding asked: while CIDone holds, a head pushed by hand without a resume is not read until the journal moves on.

## Links

### Out

- verifies [watchdog-ci-rerun](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-ci-rerun.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
