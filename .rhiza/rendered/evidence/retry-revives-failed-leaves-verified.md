# retry-revives-failed-leaves-verified

evidence · `01M6X6B152HWSYKBSCP7CK9PCH` · actor [claude-code-145f4d51](../actor/claude-code-145f4d51.md)

## result

PASS

## statement

Verifier at 01c71f0f read the diff and re-ran the acceptance. go test -trimpath -count=1 ./internal/run/... passes. In it, TestLeafFailsRetryRound has an autonomous leaf fail both rounds with its question unanswered, then Retry: the second round is Retry=true with no Decision and After=6, the builder starts at attempt 7 with spec.Retry set, the skipped dependent leaf-c builds after it passes, and the run succeeds. TestLeafFailsRetryRoundFailsAgain ends Failed after attempt 9 with arbiters still only [arbiter QL1], so there is no new arbitration. TestLeafFailsResume (answered before Retry) has exactly one leaf-round, the Decision's use-a-map, with no retry round and the Decision plus its answer in the spec. reviveLeaves gives the retry round only when no round was given in the current visit, so no leaf gets two. In internal/api, TestLeafPrompt (retry prompt names the retry and a fresh round of 3 attempts, with no Decision; answered prompt names Decision use-a-map) and TestRunInfoRetryRound (status round retry=true, decision omitempty) pass. TestScripts/describe-golden passes. The full internal/api package run on the shared VM (load ~16) failed only on environment checks: TestReadBudget timing budgets and the rootguard check of the enclosing ~/code root's state folder. No functional test failed, and the change touches neither read paths nor root state.

## Links

### Out

- verifies [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
