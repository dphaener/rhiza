# pr-246-fix-a-run-interrupted-after

evidence · `01M6K4N1JPNVZM3NZ16EEWM632` · actor [dphaener](../actor/dphaener.md)

## head

7f7f2ff6bc551c1378ad9ff82b4b15e79949b5b0

## pr

https://github.com/dphaener/rhiza/pull/246

## result

1 check passed: check SUCCESS

## statement

Pull request #246 (https://github.com/dphaener/rhiza/pull/246) passed CI at 7f7f2ff

## status

verified

## Links

### Out

- verifies [run-resume-landed-past-land](../capability/rhiza/workflows/workflow-runs/run-resume-landed-past-land.md)
- verifies [daemon-resumes-landed-past-land](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-past-land.md)
- verifies [inbox-omits-resumable-landed-run](../capability/rhiza/orchestration/inbox/inbox-omits-resumable-landed-run.md)
- verifies [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)
- verifies [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)
- verifies [run-resume](../capability/rhiza/workflows/workflow-runs/run-resume.md)
- verifies [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)
- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)
- verifies [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)
- verifies [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
