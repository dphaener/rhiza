# watchdog-pr-no-ci-verified

evidence · `01M6WVYCGDT83DTT4CR70EVND9` · actor [claude-code-bf7bb4be](../actor/claude-code-bf7bb4be.md)

## result

PASS

## statement

At 973764aa ran go test -trimpath -count=1 ./internal/api -run TestWatchdogPRNoCI -v: PASS. It uses local bare repos as fake git remotes and a fake gh PR read. In the clean case watchdogMergeMain adds one merge commit on top of the branch's old head and runs a plain git push origin HEAD:refs/heads/<branch> (no --force, no rebase). Origin's branch then equals the pushed head, and the ship state records Pushed with Head cleared. A second call does nothing. In the conflict case the merge is aborted: HEAD and origin's branch are unchanged, the worktree is clean, a.txt is reported and no ship state is written. In the pass, a clean merge resumes the run. A conflict resumes the run when the workflow has conflict handling (ship check --settled / ship rebase). Otherwise it holds a WatchdogHold that becomes an InboxStalled item. PRs that have checks or are not CONFLICTING are skipped. Read the diff. Also ran go test -run Watchdog ./internal/api (all pass), go vet and gofmt (clean). One run of the Watchdog filter tripped TestMain's guard on ~/code/.rhiza-root/usage.live.json. That file is written by the live root daemon; three reruns passed.

## Links

### Out

- verifies [watchdog-pr-no-ci](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-no-ci.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
