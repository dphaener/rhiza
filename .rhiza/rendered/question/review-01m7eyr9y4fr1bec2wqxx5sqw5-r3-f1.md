# review-01m7eyr9y4fr1bec2wqxx5sqw5-r3-f1

question · `01M7F04Q5YGDR93VPKAK43X8DB`

## statement

Plan review round 3 left finding f1 (uncovered-intent on scripts-longest-first): The intent says the scripts 'start longest first', but no leaf makes them do so. scripts-longest-first only reorders testscript's Files and says itself that nothing promises the longest scripts start first. The code bears that out: testscript's RunT (go-internal testscript.go:341-342) calls t.Run then t.Parallel for every file, so all subtests wake together when the parent returns and race for the -parallel slots in no fixed order. Locally (make test, unsharded) and within each CI shard, start order stays arbitrary, so the clause is unserved. A leaf needs a mechanism that actually orders starts, for example handing out a rank-ordered start ticket in Setup or before t.Parallel. Otherwise the intent should drop the clause and keep only the list order and the weighted shards. The critic confirmed it: The change set's intent still says cmd/rhiza's scripts 'start longest first'. After round 2's revision, scripts-longest-first's statement says it orders testscript's Files only and 'nothing here promises that the longest scripts start first', and no other leaf orders starts. script-shards-weighted only splits the scripts and orders each shard's list. The code bears this out. testscript's RunT calls t.Run, then t.Parallel first in each body (go-internal testscript.go:341-342). Each subtest then waits in tstate.waitParallel (go testing.go:1709) after the parent's barrier closes, so slots go to woken goroutines in scheduler order, not list order. The intent clause has no leaf serving it, and uncovered-intent is major, so the severity stands. How should it be settled?

## Links

### Out

- blocks [scripts-longest-first](../capability/rhiza/integrity/self-check/scripts-longest-first.md)

### In

- answers [longest-first-is-list-order](../decision/longest-first-is-list-order.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
