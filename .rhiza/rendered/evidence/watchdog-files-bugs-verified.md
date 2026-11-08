# watchdog-files-bugs-verified

evidence · `01M6WXMP7HCRZQQYJ0BWXHJ9F6` · actor [claude-code-0eca2cc4](../actor/claude-code-0eca2cc4.md)

## result

PASS

## statement

At 445732ad ran go test -trimpath -count=1 ./internal/api -run TestWatchdogBugs -v: PASS. The test drives watchdogPasses with fakes: a transient-failure recovery (index.lock) queues exactly one bug request whose text names the kind and step (transient failure at step build), the run ID, the run's request ID req1 and the log line; a second run's matching recovery queues none while the first is open (dedup by kind+step head, open = RequestState not planned); a CI flake on the wait step queues none; FileBugs=false queues none. Read the diff: fileBug is also called from the stall, unmerged-PR and PR-without-CI recoveries, the real deps append a KindBug request via daemon.UpdateRequests, and fileBugs is read from the existing rhiza.daemon/project.watchdog policy. go test -run Watchdog ./internal/api, go vet and gofmt -l on internal/api were clean.

## Links

### Out

- verifies [watchdog-files-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-files-bugs.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
