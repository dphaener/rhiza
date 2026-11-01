# answer-refuses-closed-or-unblocking

capability · `01M6AG9K13K9QP9QWR4NWM02HG` · parent [question-answer-op](../question-answer-op.md)

## acceptance

A testscript answers a blocking question, then answers it again and gets exit 2 with no answer-<q>-2 written; a second open question that blocks only another change set, answered with --change <cs>, also exits 2; a question that blocks a capability the change set touches is accepted; fold-ok passes

## statement

rhiza question answer refuses with a usage error (exit 2), writing nothing, a question that is already answered and a question that does not block the target change set, directly or through a node it touches or an ancestor of one

## Links

### Out

- decided-by [answer-reuses-blocker-rule](../../../../decision/answer-reuses-blocker-rule.md)

### In

- verifies [answer-refuses-closed-or-unblocking-verified](../../../../evidence/answer-refuses-closed-or-unblocking-verified.md)
- verifies [answer-refuses-closed-or-unblocking-qa-1](../../../../evidence/answer-refuses-closed-or-unblocking-qa-1.md)
- verifies [answer-refuses-closed-or-unblocking-verified-2](../../../../evidence/answer-refuses-closed-or-unblocking-verified-2.md)
- verifies [pr-160-question-answer-refuses-a](../../../../evidence/pr-160-question-answer-refuses-a.md)

## History

- `01M69VQG7GJARDSWZWKDJ4CWQ1` question.answer refuses a question that is not open or does not block the target change set, so an answer always lands in the change set the question blocks
