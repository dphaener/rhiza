# daemon-starts-within-max-runs

capability · `01M6XVG9DK25PFK9BHJCBDQHQ6` · parent [daemon-picks](../daemon-picks.md)

## acceptance

New testscript cmd/rhiza/testdata/script/daemon-max-runs-pending.txtar with a fake rex: with max-runs 2, one live run and a feature request whose plan run was started under 1 minute ago but has not journaled, and two ready change sets, daemon --once starts no ship run and daemon status shows both as no-slot; with max-runs 3 and the same state it starts exactly one ship run. internal/api TestDaemonScheduledStartHoldsSlotForQueuedStart: with a due schedule started and not journaled and max-runs 1, a queued start stays queued (a Go test, not a testscript: a daemon fires a schedule only on a poll after the one that first saw it, so a daemon --once process never fires one). cmd/rhiza/testdata/script/daemon-milestone-request.txtar: with --max-runs 1 and one run live, daemon --once starts nothing and the milestone and feature requests stay queued; with one interrupted run and one free slot the milestone request starts and the run is not resumed. cmd/rhiza/testdata/script/daemon-request-slot.txtar: with max-runs 2, two ready change sets and only a queued milestone request, daemon --once starts the milestone request's plan run and one ship run, never three, and daemon status shows the other as no-slot. internal/daemon round-robin/pick unit tests cover a project share computed with unjournaled starts. go test -trimpath ./internal/daemon/..., go test -trimpath ./internal/api -run TestDaemonScheduledStartHoldsSlotForQueuedStart and go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-max-runs-pending|daemon-request-slot|daemon-queued-start-slot|daemon-continue-slot|daemon-round-robin|daemon-picks|daemon-milestone-request)' pass

## statement

Every start the daemon makes in a poll (a queued milestone request's plan run, a resume of an interrupted run, a ship pick, a queued continuation, a queued start, a feature or bug request's plan run, a scheduled run) counts one machine-wide slot count against rhiza.daemon.maxRuns: runs live on the machine across every project of the daemon, plus every start in any project not yet journaled and still within daemonStartGrace (ship and queued starts held as pending, request plan runs, scheduled runs, resumes), plus the starts made earlier in the same poll; nothing starts while that count is at maxRuns. A waiting feature or bug request still keeps the next free slot (request-slot-kept, queued-start-request-slot-kept), never beyond the cap

## Links

### Out

- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- blocks [review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f1](../../../../../question/review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f1.md)
- verifies [daemon-starts-within-max-runs-verified](../../../../../evidence/daemon-starts-within-max-runs-verified.md)
- verifies [daemon-starts-within-max-runs-qa-1](../../../../../evidence/daemon-starts-within-max-runs-qa-1.md)
- verifies [daemon-starts-within-max-runs-verified-r2](../../../../../evidence/daemon-starts-within-max-runs-verified-r2.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
