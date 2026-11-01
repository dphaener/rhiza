# answer-refuses-closed-or-unblocking-verified

evidence · `01M6AZDHV51QFXYFYKFDYDFBC3` · actor [claude-code-bd7c418f](../actor/claude-code-bd7c418f.md)

## result

PASS

## statement

At f6a9bd70 I ran go test ./cmd/rhiza -run TestScripts/question-answer -count=1 (ok). The script answers store-q, then re-answers it and gets exit 2 with code usage and no answer-store-q-2; apart-q, blocking only a node the target change set does not touch, answered with --change exits 2 with no answer-apart-q; store-q blocking the touched capability store is accepted; fold-ok passes. Read the diff: runQuestionAnswer refuses via questionBlocks, which reuses planCore's blockers (open question, blocks link to the change set, a touched node or its ancestor). Probed by hand with a scratch build in a clean repo: re-answer into the same blocking change set exits 2 and writes nothing; a question blocking another accepted change set's node is refused with --change of the wrong one and accepted with the right one; a blocks link straight to the change set is accepted; a question blocking the parent of a touched child is accepted; a draft change set touching the blocked node is accepted; --record is refused outright, so the check cannot be bypassed. make check exits 0.

## Links

### Out

- verifies [answer-refuses-closed-or-unblocking](../capability/rhiza/roadmap-app-accept-answer/question-answer-op/answer-refuses-closed-or-unblocking.md)

## History

- `01M69VQG7GJARDSWZWKDJ4CWQ1` question.answer refuses a question that is not open or does not block the target change set, so an answer always lands in the change set the question blocks
