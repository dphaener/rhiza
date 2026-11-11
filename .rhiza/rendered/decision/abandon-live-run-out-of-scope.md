# abandon-live-run-out-of-scope

decision · `01M75SASV18015JAE3Q9A03FY2`

## reasoning

The change set's intent names one defect, a stale 'stalled run' inbox item, and spells out its fix; abandoning a live run is introduced only as 'also worth a look'. That second part is a different behaviour with different risk: abandon today refuses with ErrLive ('run is running in another process', internal/run/engine.go:21, internal/api/run.go:1041/1187), so making it stop a live run means signalling or taking over another process's run lock and racing the watchdog's resume, which touches the run engine's locking rather than the daemon's hold computation. Bundling it would widen a small, testable fix into a concurrency change and hold the inbox fix back. The hold fix also stands alone: once holds drop abandoned runs, the kill-plus-abandon workaround still clears the inbox. No open change set already covers stopping a live run, so it should be opened as a separate discovered-from change set.

## statement

No: this change set builds only the hold fix (watchdog holds drop runs that are abandoned or whose change set is withdrawn/rejected, pruned on the next poll). Letting rhiza run abandon (or a request withdraw) stop a live run held by another process is not a leaf here; it belongs in its own change set, linked discovered-from this one.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-32e4ca1a-1](../question/question-32e4ca1a-1.md)

## History

- `01M75S3B1YNJWFDGXCN2BEEACQ` Fix The inbox keeps a watchdog 'stalled run' item for a run that was abandoned: watchdogBudgetHolds (internal/api/daemon_watchdog.go:389) skips only runs that ended well (st.Ended && st.Result != run.Failed), and an abandoned run ends Failed, so its hold stays in daemon state and rhiza inbox lists it forever. Evidence 2026-11-11: plan run 01M75R97WXGW7XHK7T86E3470E (request 01M75R96BT1FWC9WW87NYJPHTN, change set 01M75R9CDF2TVR2CSQM99SS2P3) was stopped by the orchestrator because the owner changed the request; watchdogRequests (internal/api/watchdog_requests.go:46) saw 'its plan run is gone while the request still says planning' and resumed it twice (16:51:09Z, 16:51:43Z), the orchestrator then ran rhiza run abandon (the run is now abandoned and its change set withdrawn), yet at 17:03Z rhiza inbox still shows 'eab21e5 stalled run 01M75R97... the watchdog's retry budget of 2 is spent and it acts on this run no more'. There is no command to dismiss it. Fix: watchdogBudgetHolds (and any other hold) drops runs whose journal state is Abandoned() or whose change set is withdrawn/rejected, and existing holds for such runs are pruned on the next poll. Also worth a look: rhiza run abandon refuses while the daemon's watchdog has already resumed the run ('is running in another process'), so stopping a request's plan run takes kill-by-PID plus abandon racing the watchdog; abandon (or a rhiza request withdraw) should be able to stop a live run cleanly.
