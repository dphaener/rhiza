# review-01m6j962-r3-f4

question · `01M6JAPTEFBWESB96GAADP6AK8`

## statement

Plan review round 3 left finding f4 (contradiction on landed-run-resume-scope): The owner's decision skipped-run-in-inbox ('Show in inbox, the agent's option, picked by the owner') says an interrupted run whose change set is no longer accepted is not resumed and waits in the inbox. This assumed arbiter decision carves an exception out of that owner ruling but neither links to it nor says it refines it. The exception should name skipped-run-in-inbox and narrow it explicitly, or go to the owner as a question, so a reader can see the owner's ruling was knowingly changed. The critic confirmed it: skipped-run-in-inbox is an owner decision (wording 'Show in inbox, the agent's option, picked by the owner') stating that an interrupted run whose change set is no longer accepted is not resumed and waits in the inbox. landed-run-resume-scope, worded '(assumed, would have escalated)', resumes a subset of those runs and drops them from the inbox, yet neither its statement nor its reasoning names skipped-run-in-inbox, and it has no link to it. An arbiter decision that narrows an owner ruling without saying so contradicts it silently; it should name and narrow it explicitly (or raise it with the owner). How should it be settled?

## Links

### In

- answers [landed-run-resume-narrows-skipped-run-in-inbox](../decision/landed-run-resume-narrows-skipped-run-in-inbox.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
