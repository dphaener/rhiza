# daemon-resume-takes-slot-qa-1

evidence · `01M6YEGHAA0VF1M79Q7F9SR31A` · actor [claude-code-53686c59](../actor/claude-code-53686c59.md)

## result

PASS

## statement

QA at 3686b6ce, with the built binary in a scratch root and a fake rex/detach log. The setup had an interrupted run (a dead pid on this host) of change set one, a live run on another host and two other ready change sets. daemon --once --max-runs 1 resumed nothing and started nothing, and the journal had 0 daemon-resumed entries. --max-runs 2 printed 'Resumed run INTERRUPTED1' and no ship start; the journal had 1 daemon-resumed entry, and status showed 'runs: 1/2 live, 1 starting' with two and three 'waiting for a slot'. A second --max-runs 2 poll, with the resume still unjournaled, started nothing: the resume kept holding its slot. With the live run removed, --max-runs 3 started two ship runs and did not resume again; detach.log held one run resume and two run starts.

## Links

### Out

- verifies [daemon-resume-takes-slot](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-resume-takes-slot.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
