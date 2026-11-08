# daemon-resume-takes-slot

capability · `01M6XVG9G0WNN7JX8ZNH1C5QPG` · parent [daemon-picks](../daemon-picks.md)

## acceptance

cmd/rhiza/testdata/script/daemon-resume.txtar (or a new daemon-resume-slot.txtar) with a fake rex: with max-runs 1, one live run and one interrupted run, daemon --once resumes nothing and the run's journal has no daemon-resumed entry; with max-runs 2, one interrupted run and two ready change sets, daemon --once resumes the run and starts one ship run, not two; go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-resume' passes

## statement

The daemon resumes an interrupted run (daemon.Resumable) only while the machine-wide slot count is below rhiza.daemon.maxRuns; resumes go first among the poll's starts after a queued milestone request's plan run, each counts as a pending start for later picks in the same poll and later polls until it journals or daemonStartGrace passes, and a resume that finds no slot is left for a later poll (not journaled as resumed, so its one daemon resume is not spent)

## Links

### Out

- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [daemon-resume-takes-slot-verified](../../../../../evidence/daemon-resume-takes-slot-verified.md)
- verifies [daemon-resume-takes-slot-qa-1](../../../../../evidence/daemon-resume-takes-slot-qa-1.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
