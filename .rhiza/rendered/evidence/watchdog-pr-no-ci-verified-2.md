# watchdog-pr-no-ci-verified-2

evidence · `01M6X0J7PTDZ2GK3JYDEJSQWJJ` · actor [claude-code-b4f05f23](../actor/claude-code-b4f05f23.md)

## result

PASS

## statement

At b86358b9 ran go test -trimpath -count=1 ./internal/api -run TestWatchdog (all PASS, incl. TestWatchdogPRNoCI 1.86s), go test ./internal/daemon (ok) and go vet on both. TestWatchdogPRNoCI uses local bare repos as fake git remotes and a fake gh PR read: a clean merge of origin/main is pushed with a plain push (no --force, no rebase; origin's branch equals the pushed merge commit on top of the old head) and recorded as Pushed in the ship state; a conflicting merge is aborted with HEAD, worktree and origin's branch unchanged and no ship state written; the pass resumes a run whose workflow has conflict handling and holds an InboxStalled item for one without. Read the diff: the attempt-1 finding (merge before budget/hold checks) is fixed: watchdogPRNoCI now keeps an earlier hold for the same PR head and journal end without fetching, merging or logging (NoCIHead/NoCIAt/NoCIConflicts in daemon.WatchdogRun), and with the budget spent holds without merging, logging once; the new test cases cover both and a new head merges again.

## Links

### Out

- verifies [watchdog-pr-no-ci](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-no-ci.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
