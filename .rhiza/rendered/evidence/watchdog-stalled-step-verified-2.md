# watchdog-stalled-step-verified-2

evidence · `01M6X024KQW3AW49HE4REH97S2` · actor [claude-code-baa87db9](../actor/claude-code-baa87db9.md)

## result

PASS

## statement

Verifier at 7420691c (attempt 2): go test -trimpath -count=1 ./internal/api -run 'TestWatchdogStalled|TestWatchdogPolicy' passes. TestWatchdogStalled, with a fixed clock, a temp session store and fake send/stop/resume, nudges both autonomous runs at 21m (none at 19m), stops and resumes only run A at 41m with exactly one 'stopped and resumed' log line and Spent=1, leaves run B alone after it reported progress at 30m, and gives owner-attended run O only a hold that becomes an InboxStalled item, never a stop or resume. The fix commit's diff: StallLimit now takes the step's timeout, and a wait step gets max(waitLimit, timeout+10m), covered by the new 'wait step outlives its timeout' subtest (2h -> 2h10m, 10m -> 90m, own stall wins). run.Terminate signals -pid when the pid leads its group; detached runs start with Setsid and step commands inherit that group. go test -trimpath -count=1 ./internal/run -run TestTerminateStopsGroup passes: a group leader's background sleep child dies with it. gofmt -l and go vet are clean on internal/api and internal/run.

## Links

### Out

- verifies [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
