# longest-first-is-list-order

decision · `01M7F076QTJNMKAHYYCF7FZSX2`

## reasoning

The finding is right that list order does not fix start order: testscript's RunT calls t.Run then t.Parallel per file, and Go hands -parallel slots to woken subtests in scheduler order. But the suggested fix is unsafe: a ticket taken in Setup runs after t.Parallel has already given the script a slot, so a script holding a slot would wait for a higher-ranked script still parked in waitParallel for a slot, and with fewer slots than scripts (-parallel 8, ~400 scripts) the suite can deadlock. A safe start order needs TestScripts to schedule scripts itself instead of testscript's t.Parallel, which widens a change set the orchestrator told planners not to widen. The owner approved the milestone for its wall-time targets (cold make test within 4 minutes, CI median within 3), not for a start order; those are met by the bounded -parallel and the weighted shards, whose acceptance already measures them (shards within 15% of the mean, CI scripts jobs within 60 s of each other, -parallel 8 faster than 4). List order still helps in practice, since subtests queue for slots roughly in t.Run order. So the reviser narrows the intent's clause to what the leaves deliver, and scripts-longest-first keeps its own honest statement.

## statement

The change set does not promise that cmd/rhiza's scripts start longest first: its intent's 'start longest first' clause is narrowed to the scripts being listed longest first (scripts-longest-first) and placed longest first into weighted shards (script-shards-weighted), and no leaf adds a start-order mechanism such as a ranked start ticket

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7eyr9y4fr1bec2wqxx5sqw5-r3-f1](../question/review-01m7eyr9y4fr1bec2wqxx5sqw5-r3-f1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
