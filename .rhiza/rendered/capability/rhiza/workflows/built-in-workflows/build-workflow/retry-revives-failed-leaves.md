# retry-revives-failed-leaves

capability · `01M6WN5CRPVBSVTZ8CPGDM6AST` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test in internal/run where a leaf fails both rounds of an autonomous run with its question unanswered, then Retry: the leaf gets a leaf-round marked as the retry's, a fresh builder starts at the next attempt, its skipped dependent is built after it passes and the run succeeds; a second case where it fails again ends the run failed with no new arbitration; a third case where the leaf's own question is answered before the Retry: the leaf gets exactly one leaf-round, carrying that Decision, and no retry round; internal/api tests show run status and the builder prompt naming the retry round with no empty decision, and the answered case's builder prompt naming the Decision; go test -trimpath ./internal/run/... ./internal/api/... passes

## statement

rhiza run resume --retry on a run that failed at a leaves step gives every failed leaf of that step one fresh round of 3 attempts: a failed leaf whose own question a Decision answers by then gets that Decision's round, as on any resume, with the Decision in its builder's prompt; every other failed leaf gets a round journalled as a leaf-round given by the retry rather than by a Decision; no leaf gets two rounds from one retry; the step's skipped leaves reopen with them, a retry round's builder prompt and run status name the round as the retry's, and a failure in a retry round is not arbitrated again

## Links

### Out

- decided-by [retry-revives-every-failed-leaf](../../../../../decision/retry-revives-every-failed-leaf.md)
- derived-from [retry-revives-leaves-plan](../../../../../source/retry-revives-leaves-plan.md)

### In

- verifies [retry-revives-failed-leaves-verified](../../../../../evidence/retry-revives-failed-leaves-verified.md)
- verifies [retry-revives-failed-leaves-qa-1](../../../../../evidence/retry-revives-failed-leaves-qa-1.md)
- verifies [retry-revives-failed-leaves-verified-2](../../../../../evidence/retry-revives-failed-leaves-verified-2.md)
- verifies [retry-revives-failed-leaves-qa-2](../../../../../evidence/retry-revives-failed-leaves-qa-2.md)
- verifies [pr-377-fix-rhiza-run-resume-retry](../../../../../evidence/pr-377-fix-rhiza-run-resume-retry.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
