# plan-landed-resume-narrows-skipped-run-in-inbox

decision · `01M6K1FHTX55AC3GDV2BFJ6BRG`

## reasoning

The finding is right. skipped-run-in-inbox carries the owner's own wording and decides daemon-resumes-accepted-only and unaccepted-run-in-inbox. decision-plan-landed-resumable decides the same two nodes with a plan-run exception but never names the owner ruling, so it narrows that ruling without saying so. The change set's intent requires the exception: rhiza run resume in internal/api/run.go already lets this run through, and the inbox's advice to abandon is wrong. So the fix is to name the narrowing, not to undo it. This keeps the owner's rule in force everywhere except the one case where 'no longer accepted' was never meant to apply. The scope stays exactly what decision-plan-landed-resumable states: the plan workflow and the landed state only. That matches how the earlier review settled the same kind of silent narrowing for non-plan runs.

## statement

decision-plan-landed-resumable narrows the owner decision skipped-run-in-inbox, and says so here. skipped-run-in-inbox still holds for every interrupted run whose change set has really ended: it is not resumed and waits in the inbox with the change set's state and rhiza run abandon. The one exception is a plan-workflow run whose change set has landed. A milestone plan lands its record partway through the interview, so for a plan run landed does not mean the change set has ended, and the daemon and the inbox treat that run as in flight. A plan run whose change set was rejected, withdrawn or superseded, and every non-plan run on a change set that is no longer accepted, follow skipped-run-in-inbox unchanged.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jzdh-r3-f3](../question/review-01m6jzdh-r3-f3.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
