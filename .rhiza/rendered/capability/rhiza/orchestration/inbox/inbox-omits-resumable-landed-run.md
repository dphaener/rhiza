# inbox-omits-resumable-landed-run

capability · `01M6J9PJ6WZ9XMS1KV06HYRK4Q` · parent [inbox](../inbox.md)

## acceptance

A test in cmd/rhiza/testdata/script covers three cases with rhiza inbox --json. (1) An interrupted run on a landed change set that reached land has no skipped item. (2) The same kind of run, resumed by the daemon and then interrupted again, has a gave-up item, and one that failed a step after land has a failed item. (3) An interrupted build run on a landed change set that never reached land is still listed as skipped with rhiza run abandon

## statement

rhiza inbox does not list as skipped an interrupted non-plan run on a landed change set whose journal reached its land step, since the daemon resumes it; and it does not treat such a run's change set as gone, so if that run is interrupted again after the daemon's resume it is listed as gave-up, and if it fails a later step (such as main-green) it is listed as failed, as any other run

## Links

### Out

- depends-on [daemon-resumes-landed-past-land](../daemon/daemon-resumes-landed-past-land.md)
- decided-by [landed-run-resume-scope](../../../../decision/landed-run-resume-scope.md)

### In

- verifies [inbox-omits-resumable-landed-run-verified](../../../../evidence/inbox-omits-resumable-landed-run-verified.md)
- verifies [inbox-omits-resumable-landed-run-qa-1](../../../../evidence/inbox-omits-resumable-landed-run-qa-1.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
