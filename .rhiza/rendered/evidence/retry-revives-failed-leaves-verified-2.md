# retry-revives-failed-leaves-verified-2

evidence · `01M6X9N7NBS2PQY7YX64EK0ESC` · actor [claude-code-e805091d](../actor/claude-code-e805091d.md)

## result

PASS

## statement

Verifier, attempt 2 at a21e593c: read the fixer's diff, which handles the confirmed review findings. RunInfo.Human now prints a retry round as 'round N from attempt M: given by the retry' instead of 'decision  answers question ...', and TestRunInfoRetryRound checks that text. build-leaf-fails.txtar now expects the unanswered retry to give alpha round 2 (attempts 4-6) with no empty decision, and the answered retry to give round 3 at attempt 7 carrying use-a-map. Ran go test -trimpath ./internal/run/... (ok). The engine tests TestLeafFailsRetryRound (retry round, attempt 7 builder, then dependent leaf-c, run succeeds), TestLeafFailsRetryRoundFailsAgain (fails, arbiters only QL1) and the answered case (one leaf-round carrying the Decision, no retry) pass. Ran go test -trimpath ./internal/api/... -skip Budget: every test PASS, including TestRunInfoRetryRound, TestLeafPrompt (retry prompt with no Decision; answered prompt names the Decision) and TestRunResumeRetry. The package was flagged only by the post-run guard on ~/code/.rhiza-root/usage.live.json, which the live root daemon writes; the full api run with Budget tests ran out of time waiting on the machine-wide budget flock (load ~13), and neither is related to this change. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-leaf-fails|describe-golden)$': ok.

## Links

### Out

- verifies [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
