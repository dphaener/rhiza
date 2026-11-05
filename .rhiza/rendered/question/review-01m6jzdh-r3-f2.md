# review-01m6jzdh-r3-f2

question · `01M6K16PWTXB9T9KWBZQ1E4K9Y`

## statement

Plan review round 3 left finding f2 (contradiction on daemon-resumes-accepted-only): Overlap with accepted change set 01M6J965 (in review on cs/fix-a-run-interrupted-after, not yet on main). rhiza overlaps sees only daemon-resumes-accepted-only, but 01M6J965's branch also rewrites the statement and acceptance of unaccepted-run-in-inbox and inbox-clears-stale and the statement of daemon-resume-once, the same fields this plan rewrites, from older bases. This plan's statements also name daemon-resumes-landed-past-land and inbox-omits-resumable-landed-run, which exist only in 01M6J965, and its acceptances cite their tests. The plan has no depends-on link to 01M6J965, so it can be built or landed first, and then its nodes would name missing nodes and tests, with conflicting writes to four nodes. It needs a depends-on 01M6J965, with its updates rebased on 01M6J965's versions of those nodes. The critic confirmed it: The log ops of 01M6J965 on its branch include node.update on daemon-resumes-accepted-only (statement, acceptance), unaccepted-run-in-inbox (statement, acceptance), daemon-resume-once (statement) and inbox-clears-stale (statement, acceptance). This plan updates those same fields. Those ops are not in this worktree's log, so rhiza overlaps reports only daemon-resumes-accepted-only. daemon-resumes-landed-past-land and inbox-omits-resumable-landed-run return not_found here, yet the plan's statements name them. The plan has no depends-on link, so a depends-on 01M6J965 and a rebase onto its versions are needed. How should it be settled?

## Links

### Out

- blocks [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)

### In

- answers [r3-f2-depend-on-01m6j965](../decision/r3-f2-depend-on-01m6j965.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
