# question-32e4ca1a-1

question · `01M75S96JF35VQ1R9EP05AAENT`

## session

32e4ca1a-527d-4211-afc9-01a1cb92cf04

## statement

Change set 01M75S3 (inbox keeps a watchdog 'stalled run' item for an abandoned run): the core fix is that watchdog holds drop runs that are abandoned or whose change set is withdrawn/rejected (holds are recomputed every poll, so stale ones vanish on the next poll). The intent also suggests 'rhiza run abandon' should stop a live run cleanly instead of refusing with 'is running in another process' (today it takes kill-by-PID plus abandon racing the watchdog's resume). Should that second part be a leaf in this change set?

## Links

### In

- answers [abandon-live-run-out-of-scope](../decision/abandon-live-run-out-of-scope.md)

## History

- `01M75S3B1YNJWFDGXCN2BEEACQ` Fix The inbox keeps a watchdog 'stalled run' item for a run that was abandoned: watchdogBudgetHolds (internal/api/daemon_watchdog.go:389) skips only runs that ended well (st.Ended && st.Result != run.Failed), and an abandoned run ends Failed, so its hold stays in daemon state and rhiza inbox lists it forever. Evidence 2026-11-11: plan run 01M75R97WXGW7XHK7T86E3470E (request 01M75R96BT1FWC9WW87NYJPHTN, change set 01M75R9CDF2TVR2CSQM99SS2P3) was stopped by the orchestrator because the owner changed the request; watchdogRequests (internal/api/watchdog_requests.go:46) saw 'its plan run is gone while the request still says planning' and resumed it twice (16:51:09Z, 16:51:43Z), the orchestrator then ran rhiza run abandon (the run is now abandoned and its change set withdrawn), yet at 17:03Z rhiza inbox still shows 'eab21e5 stalled run 01M75R97... the watchdog's retry budget of 2 is spent and it acts on this run no more'. There is no command to dismiss it. Fix: watchdogBudgetHolds (and any other hold) drops runs whose journal state is Abandoned() or whose change set is withdrawn/rejected, and existing holds for such runs are pruned on the next poll. Also worth a look: rhiza run abandon refuses while the daemon's watchdog has already resumed the run ('is running in another process'), so stopping a request's plan run takes kill-by-PID plus abandon racing the watchdog; abandon (or a rhiza request withdraw) should be able to stop a live run cleanly.
