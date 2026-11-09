# daemon-starts-within-max-runs-verified

evidence · `01M6Y1GMQQB4Q0HA7Z470C7AQ2` · actor [claude-code-7afc43ea](../actor/claude-code-7afc43ea.md)

## result

PASS

## statement

Verifier at bfd78be9 (attempt 3), RHIZA_ROOT/CHANGE/RUN unset, -count=1: go test -trimpath ./internal/daemon/... ok; go test -trimpath ./internal/api -run TestDaemonScheduledStartHoldsSlotForQueuedStart PASS; go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-max-runs-pending|daemon-request-slot|daemon-queued-start-slot|daemon-continue-slot|daemon-round-robin|daemon-picks|daemon-milestone-request)' PASS, plus daemon-resume-slot, health-runs-cap and describe-golden PASS; gofmt clean, go vet ok. I read the diff d1e26ecf..bfd78be9 against the acceptance. daemon-max-runs-pending.txtar: one live run and an unjournaled feature plan run, at max-runs 2 nothing starts and status shows 2 no-slot; at 3 exactly one ship run starts. daemon-milestone-request.txtar: at --max-runs 1 with one live run nothing starts and both requests stay queued; with one interrupted run and one free slot the milestone plan run starts, no daemon-resumed entry is journaled and the feature request stays queued. daemon-request-slot.txtar (second root): with only a milestone request queued, two ready change sets and max-runs 2, it starts the plan run and one ship run (2 detach starts) and status shows 1 no-slot. TestRoundRobinCountsUnjournaledStarts covers a project share computed with unjournaled starts. daemonPoll1 now counts live runs, held, request, schedule and resume starts plus this poll's starts against max; the milestone request goes before resumes and only into a free slot. This answers attempt 1's findings f1-f3.

## Links

### Out

- verifies [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
