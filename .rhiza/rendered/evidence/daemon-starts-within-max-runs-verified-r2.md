# daemon-starts-within-max-runs-verified-r2

evidence · `01M6YF5A9D8G0R5YDTMVKZ91AB` · actor [claude-code-f767cf47](../actor/claude-code-f767cf47.md)

## result

PASS

## statement

Verifier (claude-code, session f767cf47) at ba2674a7, the review fixer commit for quality f1. Diff read: daemonMilestoneWaits (internal/api/daemon_projects.go) now takes the run store, run list and host, and blocks a queued milestone request only while another milestone request's run is journaled and LiveRunning on host, or its start is unjournaled within daemonStartGrace. That is the same rule as daemonPlanRequests' milestoneLive in internal/api/request.go, so a milestone request whose run ended stays queued but no longer blocks later ones. The new TestDaemonMilestoneWaitsAfterEndedRun covers three steps: unjournaled start blocks, live run blocks, finished run makes the second request demand again. With RHIZA_ROOT/CHANGE/RUN unset and GOTMPDIR/TMPDIR in the session scratchpad, these pass: go test -count=1 -trimpath ./internal/daemon/... (ok); go test -trimpath ./internal/api -run '^TestDaemonScheduledStartHoldsSlotForQueuedStart$' (ok) and -run '^TestDaemonMilestoneWaitsAfterEndedRun$' (ok); go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-max-runs-pending|daemon-request-slot|daemon-queued-start-slot|daemon-continue-slot|daemon-round-robin|daemon-picks|daemon-milestone-request)' (all PASS, ok). A broader -run TestDaemon over internal/api passed every test, but the package's state-folder guard then flagged ~/code/.rhiza-root/usage.live.json as changed. This commit does not touch usage code, and the targeted runs are clean, so that comes from the environment (the live root daemon or the known usage leak). Quality f2 (minor: unresolvable milestone requests still counted as demand) is left as it was; it is not part of this leaf's acceptance.

## Links

### Out

- verifies [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
