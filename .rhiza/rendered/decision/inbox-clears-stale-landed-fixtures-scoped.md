# inbox-clears-stale-landed-fixtures-scoped

decision · `01M6K1K7A51YFGR1S5PZK8Y1DC`

## reasoning

The capability's statement already excludes two kinds of run on a landed change set from being left out: a plan-workflow run (inbox-landed-plan-run-in-flight) and a non-plan run whose journal reached land (inbox-omits-resumable-landed-run). The acceptance's first clause says 'a failed run and a gave-up run each on a landed ... change set' are omitted with no run kind, so 'exactly' contradicts the exceptions listed after it, as the plan reviewer and critic both confirmed. Scoping the landed fixtures to non-plan runs that never reached land is the only reading consistent with the statement and with this change set's intent (treat an interrupted plan run whose record landed as resumable); it changes no behaviour, only makes the check unambiguous, and it is the wording 01M6J965 already uses.

## statement

Reword inbox-clears-stale's acceptance so the 'lists exactly the accepted-change-set items' clause names its landed fixtures as a failed and a gave-up non-plan run whose journal never reached land (the withdrawn, rejected and superseded fixtures stay any run); the plan-run and reached-land exceptions then add to that set rather than contradict it, matching 01M6J965's version of the acceptance

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jzdh-r3-f4](../question/review-01m6jzdh-r3-f4.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
