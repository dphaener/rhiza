# review-01m6j962-r3-f3

question · `01M6JAPK0HXVDH2NVD6H2G15B8`

## statement

Plan review round 3 left finding f3 (contradiction on daemon-resumes-accepted-only): The plan rewrites the statement to add the landed-past-land exception, but the acceptance still reads 'leaves interrupted runs on landed, rejected, withdrawn, proposed and draft change sets unresumed'. Unqualified, that contradicts the new statement. The acceptance needs the same split: landed and never reached land means unresumed, landed and reached land means resumed. The critic confirmed it: The plan rewrote daemon-resumes-accepted-only's statement to add the landed-past-land exception but left its acceptance unqualified: 'leaves interrupted runs on landed ... change sets unresumed'. The existing fixtures never reach land, so the test still passes, but the acceptance and statement of the same node now disagree on landed runs, and a builder reading the acceptance could take it as requiring every landed run to stay unresumed. It should say landed runs that never reached land. How should it be settled?

## Links

### Out

- blocks [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)

### In

- answers [daemon-resumes-accepted-only-acceptance-split](../decision/daemon-resumes-accepted-only-acceptance-split.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
