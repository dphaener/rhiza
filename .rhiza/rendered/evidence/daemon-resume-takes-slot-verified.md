# daemon-resume-takes-slot-verified

evidence · `01M6Y1FA81WN5NZ3C6PEJ7GPAR` · actor [claude-code-7fc5c679](../actor/claude-code-7fc5c679.md)

## result

PASS

## statement

At 685fd32d, cmd/rhiza/testdata/script/daemon-resume-slot.txtar with a fake rex checks both acceptance cases. Case 1: --max-runs 1, one live run (on another host) and one interrupted run. daemon --once prints no 'Resumed run' and no 'Started a run of', writes no detach.log, and INTERRUPTED1.jsonl has 0 daemon-resumed entries. Case 2: the live run is removed, then --max-runs 2 with interrupted INTERRUPTED1 (change set One) and two other ready change sets. It prints 'Resumed run INTERRUPTED1', exactly one 'Started a run of' line (stdout -count=1), one daemon-resumed journal entry, and one 'run resume' plus one 'run start' in detach.log. I ran go test -count=1 -trimpath ./cmd/rhiza -run 'TestScripts/daemon-resume' with RHIZA_ROOT/CHANGE/RUN unset and GOTMPDIR/TMPDIR in the session scratchpad: daemon-resume, daemon-resume-slot, daemon-resume-landed and daemon-resume-landed-plan all PASS. go test -trimpath -run 'Daemon|Resum|Hold' ./internal/api printed PASS. Rootguard then flagged ~/code/.rhiza-root/usage.live.json, which the live root daemon writes; that is the known usage leak, not this code. From the diff 5ba236d4..685fd32d of internal/api/daemon.go and daemon_projects.go: a queued milestone request's plan run is picked before the resumes. Each resume is skipped while used >= max, and that check runs before the daemon-resumed entry is journaled, so a resume with no slot leaves its one resume unspent. Each resume increments used and counts in the ship picks' extra. The new daemonResumesStarting keeps counting a resumed run in later polls (and in daemonDemand) until its process journals or daemonStartGrace passes.

## Links

### Out

- verifies [daemon-resume-takes-slot](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-resume-takes-slot.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
