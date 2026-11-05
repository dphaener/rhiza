# daemon-resumes-landed-past-land

capability · `01M6J9P6J78MZ74PAM0VB6ZD5W` · parent [daemon](../daemon.md)

## acceptance

A unit test of daemon.Resumable returns an interrupted run on a landed change set whose journal entered land, and not one on a landed change set whose journal never reached land, nor runs on rejected, withdrawn, proposed or draft change sets; daemon-resume.txtar gains a run interrupted after land on a landed change set that gets one daemon-resumed entry and whose journal then shows the run going on past its stopped step to succeed, with no 'no longer accepted' refusal

## statement

The daemon resumes, once and headless, an interrupted non-plan run whose change set has landed when its journal shows the run reached its land step, as it resumes a run on an accepted change set

## Links

### Out

- depends-on [daemon-resume-once](daemon-resume-once.md)
- depends-on [run-resume-landed-past-land](../../workflows/workflow-runs/run-resume-landed-past-land.md)
- decided-by [landed-run-resume-scope](../../../../decision/landed-run-resume-scope.md)

### In

- depends-on [inbox-omits-resumable-landed-run](../inbox/inbox-omits-resumable-landed-run.md)
- verifies [daemon-resumes-landed-past-land-verified](../../../../evidence/daemon-resumes-landed-past-land-verified.md)
- verifies [daemon-resumes-landed-past-land-qa-1](../../../../evidence/daemon-resumes-landed-past-land-qa-1.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
