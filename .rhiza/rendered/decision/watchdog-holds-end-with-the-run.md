# watchdog-holds-end-with-the-run

decision · `01M75SBSGPAWT9WZF7YE3HENX5`

## reasoning

The change set's intent prescribes this fix. Abandoning a run is the owner's (or orchestrator's) explicit dismissal, and a withdrawn or rejected change set has nothing left to recover, so a hold on either is an inbox item nobody can act on and no command clears. Recomputing holds from journals and change-set states each poll needs no migration or prune step for holds already saved.

## statement

A watchdog hold lasts only while its run can still be acted on: a run that is abandoned, or whose change set is withdrawn or rejected, gets no hold, and since holds are recomputed each poll its existing hold is gone at the next poll

## Links

### Out

- derived-from [abandoned-run-hold-plan](../source/abandoned-run-hold-plan.md)

### In

- decided-by [watchdog-holds-skip-abandoned](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-holds-skip-abandoned.md)

## History

- `01M75S3B1YNJWFDGXCN2BEEACQ` Fix The inbox keeps a watchdog 'stalled run' item for a run that was abandoned: watchdogBudgetHolds (internal/api/daemon_watchdog.go:389) skips only runs that ended well (st.Ended && st.Result != run.Failed), and an abandoned run ends Failed, so its hold stays in daemon state and rhiza inbox lists it forever. Evidence 2026-11-11: plan run 01M75R97WXGW7XHK7T86E3470E (request 01M75R96BT1FWC9WW87NYJPHTN, change set 01M75R9CDF2TVR2CSQM99SS2P3) was stopped by the orchestrator because the owner changed the request; watchdogRequests (internal/api/watchdog_requests.go:46) saw 'its plan run is gone while the request still says planning' and resumed it twice (16:51:09Z, 16:51:43Z), the orchestrator then ran rhiza run abandon (the run is now abandoned and its change set withdrawn), yet at 17:03Z rhiza inbox still shows 'eab21e5 stalled run 01M75R97... the watchdog's retry budget of 2 is spent and it acts on this run no more'. There is no command to dismiss it. Fix: watchdogBudgetHolds (and any other hold) drops runs whose journal state is Abandoned() or whose change set is withdrawn/rejected, and existing holds for such runs are pruned on the next poll. Also worth a look: rhiza run abandon refuses while the daemon's watchdog has already resumed the run ('is running in another process'), so stopping a request's plan run takes kill-by-PID plus abandon racing the watchdog; abandon (or a rhiza request withdraw) should be able to stop a live run cleanly.
