# review-01m6j962-r3-f2

question · `01M6JAPAZQ4ME30QT47YRE8CZ7`

## statement

Plan review round 3 left finding f2 (contradiction on daemon-resume-once): Its statement still says the daemon 'leaves alone an interrupted run whose change set is not accepted'. daemon-resumes-landed-past-land has the daemon resume an interrupted run on a landed change set. The plan should narrow this statement as it did for daemon-resumes-accepted-only, with the same exception. The critic confirmed it: daemon-resume-once (version 01M6GBJX, not changed by this change set) still says the daemon 'leaves alone an interrupted run whose change set is not accepted'. daemon-resumes-landed-past-land has the daemon resume such a run on a landed change set. The statement needs the same exception the plan gave daemon-resumes-accepted-only. How should it be settled?

## Links

### Out

- blocks [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

### In

- answers [narrow-daemon-resume-once](../decision/narrow-daemon-resume-once.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
