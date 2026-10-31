# question-answer-op

capability · `01M686TPQBCC7169KT6WEPB85A` · parent [roadmap-app-accept-answer](../roadmap-app-accept-answer.md)

## acceptance

A testscript where a question blocks an accepted change set: rhiza question answer <q> <words> --change <cs> --json makes ready drop the reason and fold-ok passes; --dry-run writes nothing; answering an unknown question exits 2; two answers to two questions get distinct handles; the roadmap app's answer path calls question.answer (Go test with a fake runner); describe-golden.txtar is regenerated

## statement

A question.answer operation answers an open question that blocks a change set by recording the owner's words as a Decision quoting them, linked answers to the question, under a free handle, in the change set it blocks; the roadmap app answers through it

## Links

### Out

- decided-by [inbox-blockers-per-change-set](../../../decision/inbox-blockers-per-change-set.md)

### In

- verifies [question-answer-op-verified](../../../evidence/question-answer-op-verified.md)
- verifies [question-answer-op-qa-1](../../../evidence/question-answer-op-qa-1.md)
- verifies [question-answer-op-verified-2](../../../evidence/question-answer-op-verified-2.md)
- verifies [pr-143-inbox-shows-blockers](../../../evidence/pr-143-inbox-shows-blockers.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
