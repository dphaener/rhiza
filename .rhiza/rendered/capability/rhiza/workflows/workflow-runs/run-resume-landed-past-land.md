# run-resume-landed-past-land

capability · `01M6J9NTD1SF9Z8B72JJ0MV1HG` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript resumes a ship-change-shaped run interrupted after its land step on a landed change set and it goes on past its stopped step to succeed; in the same suite run-resume-landed-plan.txtar keeps passing (its ship-change run BUILDLAND01 on a landed change set whose journal never reached land is refused with exit 3 naming rhiza run abandon, journal unchanged) and run-resume-unaccepted.txtar keeps passing (its run on a withdrawn change set is refused the same way, journal unchanged)

## statement

rhiza run resume continues an unfinished non-plan run whose change set has landed when its journal shows the run reached its land step (land entered or passed), going on from where the journal left it; a non-plan run on a landed change set that never reached land is still refused as no longer accepted, naming rhiza run abandon

## Links

### Out

- depends-on [run-resume](run-resume.md)
- decided-by [landed-run-resume-scope](../../../../decision/landed-run-resume-scope.md)

### In

- depends-on [daemon-resumes-landed-past-land](../../orchestration/daemon/daemon-resumes-landed-past-land.md)
- verifies [run-resume-landed-past-land-verified](../../../../evidence/run-resume-landed-past-land-verified.md)
- verifies [run-resume-landed-past-land-qa-1](../../../../evidence/run-resume-landed-past-land-qa-1.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
