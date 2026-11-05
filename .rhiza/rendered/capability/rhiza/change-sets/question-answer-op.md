# question-answer-op

capability · `01M686TPQBCC7169KT6WEPB85A` · parent [change-sets](../change-sets.md)

## acceptance

A testscript where a question blocks an accepted change set: rhiza question answer <q> <words> --change <cs> --json makes ready drop the reason and fold-ok passes; --dry-run writes nothing; answering an unknown question exits 2; two answers to two questions get distinct handles; describe-golden.txtar is regenerated

## statement

A question.answer operation answers an open question that blocks a change set by recording the owner's words as a Decision quoting them, linked answers to the question, under a free handle, in the change set it blocks

## Links

### Out

- decided-by [inbox-blockers-per-change-set](../../../decision/inbox-blockers-per-change-set.md)

### In

- verifies [question-answer-op-verified](../../../evidence/question-answer-op-verified.md)
- verifies [question-answer-op-qa-1](../../../evidence/question-answer-op-qa-1.md)
- verifies [question-answer-op-verified-2](../../../evidence/question-answer-op-verified-2.md)
- verifies [pr-143-inbox-shows-blockers](../../../evidence/pr-143-inbox-shows-blockers.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
