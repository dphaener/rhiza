# watchdog-poll-behind-verified

evidence · `01M6WY90T2YQC7GTZV7Y01SBKP` · actor [claude-code-09c3c21a](../actor/claude-code-09c3c21a.md)

## result

PASS

## statement

At 315c91c7, ran go test -trimpath -count=1 ./internal/api -run TestWatchdogPollBehind -v: PASS. The test uses fixed fake times and durations: a poll at the 3x5s limit (15s) logs nothing and health stays ok; a 40s poll logs exactly one line naming 40s and 15s and turns the daemon health part to warn with 'polling behind'; a 2s poll clears state.PollBehind and health is ok again with no new log line; a disabled watchdog records nothing. Read the diff: runDaemon times each poll with ctx.Now, applies the home project's watchdog pollBehind policy (default 3) to its *daemonProj state before saveState, carries PollBehind across restarts, daemon status reports poll_behind and healthDaemon warns from it; describe-golden updated. Also ran TestDaemon*/TestHealth* in ./internal/api (all pass; the package's known enclosing-root state leak comes from existing tests, the new test alone is clean), ./internal/daemon tests, and cmd/rhiza TestScripts for describe-golden, health* and daemon* (ok); go vet and gofmt clean.

## Links

### Out

- verifies [watchdog-poll-behind](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-poll-behind.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
