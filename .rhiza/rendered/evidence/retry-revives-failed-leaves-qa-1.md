# retry-revives-failed-leaves-qa-1

evidence · `01M6X6YPEMKPRCMHM0EWX7SQ0X` · actor [claude-code-5a9b12d9](../actor/claude-code-5a9b12d9.md)

## result

FAIL

## statement

QA with the branch's built rhiza binary (fake rex/claude, scratch root): leaf alpha failed 3 attempts (beta skipped); run resume --retry with its question unanswered gave alpha a retry round (JSON round 2 retry:true, no decision) and attempts 4-6 ran; with the builder then passing, a second retry gave round 3, the builder prompt named the retry round with no Decision, alpha passed at attempt 7, skipped beta was built and the run succeeded. In a fresh run with the question answered by Decision use-a-map before the retry, alpha got exactly one round (decision use-a-map, no retry round, no 'retry gives it' log line) and its prompt named the Decision. Problems: human run status prints 'decision  answers question' (empty) for retry rounds (internal/api/run.go:2952), and testscript build-leaf-fails.txtar still asserts the old 'unanswered retry starts nothing' and fails.

## Links

### Out

- verifies [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
