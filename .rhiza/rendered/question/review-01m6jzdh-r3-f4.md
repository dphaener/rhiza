# review-01m6jzdh-r3-f4

question · `01M6K1767YD1SEHPSD4MM23BKT`

## statement

Plan review round 3 left finding f4 (contradiction on inbox-clears-stale): The new acceptance still reads 'a failed run and a gave-up run each on a landed ... change set ... lists exactly the accepted-change-set items', then adds that failed and gave-up plan runs on a landed change set stay listed. The first clause doesn't say which runs it means, so it contradicts the exceptions that follow. It should say the landed fixtures are non-plan runs whose journal never reached land, as 01M6J965's version of this acceptance does. The critic confirmed it: inbox-clears-stale's acceptance in this plan reads: a failed and a gave-up run on a landed (and withdrawn, rejected, superseded) change set 'lists exactly the accepted-change-set items'. It then says failed and gave-up plan runs on a landed change set, and non-plan runs that reached land, stay listed. The first clause doesn't say which kind of landed fixtures it means, so 'exactly' contradicts the exceptions after it. It should say those landed fixtures are non-plan runs whose journal never reached land. How should it be settled?

## Links

### Out

- blocks [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

### In

- answers [inbox-clears-stale-landed-fixtures-scoped](../decision/inbox-clears-stale-landed-fixtures-scoped.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
