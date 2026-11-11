# watchdog-holds-skip-abandoned-qa-1

evidence · `01M76CYJYJJPHDD90J95CGN83F` · actor [claude-code-9ba60f4d](../actor/claude-code-9ba60f4d.md)

## result

PASS

## statement

QA at dc287cd2d with the built binary: scratch Rhiza root (own HOME) with one registered project, four autonomous failed runs whose change sets are draft, draft, withdrawn (rhiza withdraw) and rejected (rhiza propose + reject), and a daemon state seeding each with a spent budget (2 of 2) and a saved hold; rhiza inbox listed 4 'stalled run' items. Ran rhiza run abandon on the second run (inbox still 4 until a poll, as specified), then rhiza daemon --once --max-runs 0: state watchdog_holds kept only the plain failed run, rhiza inbox showed 1 stalled item (that run, with its tried attempts), budget memory for all four kept; a second poll stayed at 1. Control: re-seeding the 4 holds and polling with a build of the merge-base kept all 4 stalled items; the change's build then dropped them back to 1. go test -trimpath ./internal/api -run 'TestWatchdog(Budget|Holds)' PASS.

## Links

### Out

- verifies [watchdog-holds-skip-abandoned](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-holds-skip-abandoned.md)

## History

- `01M75S3B1YNJWFDGXCN2BEEACQ` Fix The inbox keeps a watchdog 'stalled run' item for a run that was abandoned: watchdogBudgetHolds (internal/api/daemon_watchdog.go:389) skips only runs that ended well (st.Ended && st.Result != run.Failed), and an abandoned run ends Failed, so its hold stays in daemon state and rhiza inbox lists it forever. Evidence 2026-11-11: plan run 01M75R97WXGW7XHK7T86E3470E (request 01M75R96BT1FWC9WW87NYJPHTN, change set 01M75R9CDF2TVR2CSQM99SS2P3) was stopped by the orchestrator because the owner changed the request; watchdogRequests (internal/api/watchdog_requests.go:46) saw 'its plan run is gone while the request still says planning' and resumed it twice (16:51:09Z, 16:51:43Z), the orchestrator then ran rhiza run abandon (the run is now abandoned and its change set withdrawn), yet at 17:03Z rhiza inbox still shows 'eab21e5 stalled run 01M75R97... the watchdog's retry budget of 2 is spent and it acts on this run no more'. There is no command to dismiss it. Fix: watchdogBudgetHolds (and any other hold) drops runs whose journal state is Abandoned() or whose change set is withdrawn/rejected, and existing holds for such runs are pruned on the next poll. Also worth a look: rhiza run abandon refuses while the daemon's watchdog has already resumed the run ('is running in another process'), so stopping a request's plan run takes kill-by-PID plus abandon racing the watchdog; abandon (or a rhiza request withdraw) should be able to stop a live run cleanly.
