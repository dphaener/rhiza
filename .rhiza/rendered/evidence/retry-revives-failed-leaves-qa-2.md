# retry-revives-failed-leaves-qa-2

evidence · `01M6X9V6YKB4M1HV45FVD6EZVG` · actor [claude-code-c5af36ea](../actor/claude-code-c5af36ea.md)

## result

PASS

## statement

QA cycle 2 with the branch's built rhiza binary at b7586ebf (scratch root, fake rex/claude, isolated HOME): leaf alpha failed 3 attempts, beta skipped, gamma passed. run resume --retry with the question unanswered gave alpha a retry round: JSON round 2 retry:true with no decision, attempts 4-6 ran, a new question leaf-failed-...-2 blocks it, and human run status now prints 'round 2 from attempt 4: given by the retry' (no empty 'decision  answers'). The retry-round builder prompts say 'A person retried the run, which gives it a fresh round of 3 attempts; no Decision answers its failure' (round-1 prompts do not). A second retry with the builder passing gave round 3 (also the retry's), alpha passed at attempt 7, skipped beta was built and the run succeeded. In a fresh run with Decision use-a-map answering the leaf's question before the retry, alpha got exactly one round (round 2 from attempt 4: decision use-a-map answers question ...), the engine logged 'gets round 2: decision use-a-map', no leaf round had retry:true, the builder prompt named the Decision and not the retry, and the run succeeded. go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-leaf-fails|describe-golden)$', ./internal/run -run TestLeafFails and ./internal/api -run 'TestLeafPrompt|TestRunInfoRetryRound|TestRunResumeRetry' pass.

## Links

### Out

- verifies [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
