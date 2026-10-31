# inbox-blocking-questions-qa-1

evidence · `01M69VVSECRMZDH7CBVN971383` · actor [claude-code-77fbcb0f](../actor/claude-code-77fbcb0f.md)

## result

PASS

## statement

At 302d9b36, built rhiza and in a scratch repo made a landed change set and a draft each with a recorded open question blocking a node they touch: rhiza inbox printed Nothing waits on you. Then an accepted change set touching store with store-q blocking store: rhiza inbox showed a blocker item (question: Which store? (store-q), blocks: store, change set Accepted) and --json kind blocker with handle/blocks/question. Adding a run-log question on the same change set and a proposed change set blocked directly by cs-q listed the run question, the store-q blocker, and a blocks: this change set blocker together. After answering store-q the blocker was gone (Nothing waits on you).

## Links

### Out

- verifies [inbox-blocking-questions](../capability/rhiza/orchestration/inbox/inbox-blocking-questions.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
