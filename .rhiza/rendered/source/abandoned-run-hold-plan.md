# abandoned-run-hold-plan

source · `01M75SBF6DHCZDGE9PVRVJSS99`

## quotes

```json
[]
```

## raw_sha256

cc3194d1bdaad6df2412ade95e4e7e1aa5b995d71d15815a898659c341df1e83

## session

32e4ca1a-527d-4211-afc9-01a1cb92cf04

## summary

Autonomous plan run on change set 01M75S3 (no owner present). The planner read internal/api/daemon_watchdog.go: watchdogBudgetHolds skips only runs that ended well, so an abandoned run (ended Failed) keeps its 'retry budget spent' hold, and rhiza inbox lists it with no way to dismiss it. Holds are recomputed every poll (state.WatchdogHolds = holds in the watchdog poll), so skipping abandoned runs, and runs whose change set is withdrawn or rejected, also prunes existing holds on the next poll. The planner asked whether letting rhiza run abandon stop a live run (today refused with 'is running in another process') belongs here; the arbiter ruled no (Decision abandon-live-run-out-of-scope): only the hold fix is built here, and stopping a live run goes to its own change set, linked discovered-from this one. The change set serves m4-runs-hold-up, like the earlier watchdog fix 01M73G0W.

## title

Plan of the abandoned-run watchdog hold fix (autonomous)

## Links

### In

- derived-from [watchdog-holds-end-with-the-run](../decision/watchdog-holds-end-with-the-run.md)
- derived-from [watchdog-holds-skip-abandoned](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-holds-skip-abandoned.md)

## History

- `01M75S3B1YNJWFDGXCN2BEEACQ` Fix The inbox keeps a watchdog 'stalled run' item for a run that was abandoned: watchdogBudgetHolds (internal/api/daemon_watchdog.go:389) skips only runs that ended well (st.Ended && st.Result != run.Failed), and an abandoned run ends Failed, so its hold stays in daemon state and rhiza inbox lists it forever. Evidence 2026-11-11: plan run 01M75R97WXGW7XHK7T86E3470E (request 01M75R96BT1FWC9WW87NYJPHTN, change set 01M75R9CDF2TVR2CSQM99SS2P3) was stopped by the orchestrator because the owner changed the request; watchdogRequests (internal/api/watchdog_requests.go:46) saw 'its plan run is gone while the request still says planning' and resumed it twice (16:51:09Z, 16:51:43Z), the orchestrator then ran rhiza run abandon (the run is now abandoned and its change set withdrawn), yet at 17:03Z rhiza inbox still shows 'eab21e5 stalled run 01M75R97... the watchdog's retry budget of 2 is spent and it acts on this run no more'. There is no command to dismiss it. Fix: watchdogBudgetHolds (and any other hold) drops runs whose journal state is Abandoned() or whose change set is withdrawn/rejected, and existing holds for such runs are pruned on the next poll. Also worth a look: rhiza run abandon refuses while the daemon's watchdog has already resumed the run ('is running in another process'), so stopping a request's plan run takes kill-by-PID plus abandon racing the watchdog; abandon (or a rhiza request withdraw) should be able to stop a live run cleanly.
