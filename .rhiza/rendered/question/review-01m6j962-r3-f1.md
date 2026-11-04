# review-01m6j962-r3-f1

question · `01M6JAP2M7P69YCNC18EB7XT3Y`

## statement

Plan review round 3 left finding f1 (contradiction on unaccepted-run-in-inbox): This verified capability says an interrupted run on a change set no longer accepted waits in the inbox, and its acceptance tests that an interrupted run on a landed change set is listed with rhiza run abandon and that the daemon does not resume it. The plan's inbox-omits-resumable-landed-run and daemon-resumes-landed-past-land do the opposite for a non-plan run that reached land. The plan narrows run-resume, manual-resume-refuses-unaccepted and daemon-resumes-accepted-only, but leaves this capability's statement and acceptance as they are, so the model contradicts itself and the existing test's claim gets weaker without anyone noting it. The critic confirmed it: unaccepted-run-in-inbox's statement says any interrupted run on a change set no longer accepted waits in the inbox. The plan's inbox-omits-resumable-landed-run removes the skipped item for a non-plan run on a landed change set that reached land, and the plan does not touch this node, while it did narrow the sibling nodes run-resume, manual-resume-refuses-unaccepted and daemon-resumes-accepted-only. Its existing test (daemon-resume.txtar LANDED00001/LANDEDAUTO1, journals with no step entries) still passes, so only the statement and acceptance need narrowing to runs that never reached land, but as the model stands it contradicts itself. How should it be settled?

## Links

### Out

- blocks [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)

### In

- answers [narrow-unaccepted-run-in-inbox](../decision/narrow-unaccepted-run-in-inbox.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
