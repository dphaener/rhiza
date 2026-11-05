# r3-f2-depend-on-01m6j965

decision · `01M6K1CMDJKW4ACQ2S6BVH6HSV`

## reasoning

The finding was confirmed by the critic: 01M6J965 is accepted and in review on cs/fix-a-run-interrupted-after, not on main, and its branch log updates the same statement and acceptance fields on four nodes this plan rewrites, and adds two nodes this plan names. Without ordering, this plan could build or land first, naming missing nodes and tests and making conflicting writes from older bases. Both change sets fix the same area (resuming runs whose change set landed), so 01M6J965 is the foundation and this one refines it; dropping the references or folding the two together would lose either plan's scope. A depends-on link plus rebasing onto 01M6J965's field versions is the smallest fix and exactly what the finding asks for.

## statement

This change set depends-on 01M6J965BRJ2ZTDT7Y311PM96G and is built and landed only after it. Its updates to daemon-resumes-accepted-only, unaccepted-run-in-inbox, inbox-clears-stale and daemon-resume-once are rebased onto 01M6J965's versions of those fields (keeping 01M6J965's wording and adding only what this fix needs), and its references to daemon-resumes-landed-past-land and inbox-omits-resumable-landed-run and their tests stay, resolved once 01M6J965 lands.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jzdh-r3-f2](../question/review-01m6jzdh-r3-f2.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
