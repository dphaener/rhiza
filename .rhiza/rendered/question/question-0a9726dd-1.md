# question-0a9726dd-1

question · `01M6J9DNGAX992WSTAEVHN2HPY`

## session

0a9726dd-c836-4069-833f-01a1a4935119

## statement

Change set 01M6J965 (resume runs stranded after their change set landed) serves no intent yet. Which landed intent should it serve? m4-runs-hold-up says 'A build, review or ship-change run finishes without hand repair', which is what this fixes. / Draft change set 01M6H5VQ already covers the plan-run half: inbox and daemon treat an interrupted plan run whose record landed as resumable (rhiza run resume already continues one, plan-run-resumes-landed-record). Should 01M6J965 leave plan runs to it and cover only non-plan (ship-change) runs? / What makes a ship-change run on a landed change set resumable? `rhiza ship land` already skips a change set that has landed, so re-running it is safe.

## Links

### In

- answers [landed-run-resume-scope](../decision/landed-run-resume-scope.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
