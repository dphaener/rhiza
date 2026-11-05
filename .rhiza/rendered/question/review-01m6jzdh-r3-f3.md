# review-01m6jzdh-r3-f3

question · `01M6K1707S7Z6354MXB7R18GYF`

## statement

Plan review round 3 left finding f3 (contradiction on decision-plan-landed-resumable): The owner decision skipped-run-in-inbox ('Show in inbox, the agent's option, picked by the owner') says an interrupted run whose change set is no longer accepted is not resumed and waits in the inbox. This assumed arbiter decision makes an exception for plan runs on a landed change set and is linked as decided-by to unaccepted-run-in-inbox and daemon-resumes-accepted-only, both of which skipped-run-in-inbox also decides. But it never names or links skipped-run-in-inbox. That is the same silent narrowing of an owner ruling that 01M6J965's review flagged (review-01m6j962-r3-f4) and settled with landed-run-resume-narrows-skipped-run-in-inbox, which covers only non-plan runs. The decision should name skipped-run-in-inbox and narrow it explicitly. The critic confirmed it: skipped-run-in-inbox carries the owner's wording ('Show in inbox, the agent's option, picked by the owner'): an interrupted run on a change set no longer accepted is not resumed and waits in the inbox. decision-plan-landed-resumable (wording '(assumed, would have escalated)') makes plan runs on landed change sets an exception to that rule. Yet no op in this change set's log names or links skipped-run-in-inbox; it links only answers and derived-from. The owner ruling is narrowed without saying so. It should be named and narrowed explicitly, as 01M6J965 did for non-plan runs. How should it be settled?

## Links

### In

- answers [plan-landed-resume-narrows-skipped-run-in-inbox](../decision/plan-landed-resume-narrows-skipped-run-in-inbox.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
