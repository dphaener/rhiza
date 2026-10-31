# inbox-blocking-questions

capability · `01M686TR0KYVRWTYHXJ3F8P1D0` · parent [inbox](../inbox.md)

## acceptance

inbox-blocking-questions.txtar with one run question and one model question blocking an accepted change set shows both in text and --json; after the model question is answered it is gone; a question on a draft or landed change set is never listed

## statement

rhiza inbox lists each open question that blocks a proposed or accepted change set as a blocker item, giving the question, the change set and what it blocks, next to the run items

## Links

### Out

- decided-by [inbox-blockers-per-change-set](../../../../decision/inbox-blockers-per-change-set.md)

### In

- verifies [inbox-blocking-questions-verified](../../../../evidence/inbox-blocking-questions-verified.md)
- verifies [inbox-blocking-questions-qa-1](../../../../evidence/inbox-blocking-questions-qa-1.md)
- verifies [pr-143-inbox-shows-blockers](../../../../evidence/pr-143-inbox-shows-blockers.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
