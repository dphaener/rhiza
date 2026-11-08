# watchdog-stalled-step-verified

evidence · `01M6WRNP2EXMHAPW10JMK1NC98` · actor [claude-code-e07461be](../actor/claude-code-e07461be.md)

## result

PASS

## statement

Verifier claude-code-e07461be at b0c7d417: ran go test -trimpath -count=1 ./internal/api -run TestWatchdogStalled -v: PASS. The test uses a fixed clock (t0 plus offsets), a temp session store with fake hook reports and fake send/stop/resume/log deps. Autonomous run A: no action at 19m, nudged at 21m (first limit, agentNudge 20m), stopped and resumed once at 41m (second limit, agentRetry 40m) with exactly one 'stopped and resumed' log line and Spent=1; at 90m, with retryBudget 1 spent, it is held instead of retried. Run B reported progress at 30m, so it was left alone at 41m. Owner-attended run O was never stopped or resumed, got a log line and a hold, and the hold became an InboxStalled item. Read the diff: watchdogStalled/daemonWatchdog in internal/api/daemon_watchdog.go, wired into the poll in place of daemonNudge when watchdog.enabled; memory and holds persisted in daemon.State and carried across restarts; holds surfaced in inbox and notify; session.Stallable split out of Stalled with the same behaviour; run.Terminate sends SIGTERM. Also passed: TestWatchdogPolicy and the internal/session, run and daemon packages; go vet is clean. All tests in the wider api subset (Nudge|Inbox|Daemon|Stall) passed, but the package's root-state guard then flagged ~/code/.rhiza-root/usage.live.json, which the live root daemon on the VM writes during a 575s run. That is not touched by this diff.

## Links

### Out

- verifies [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
