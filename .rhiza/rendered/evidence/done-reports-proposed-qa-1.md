# done-reports-proposed-qa-1

evidence · `01M658X8HRQANQK6Q43V19986C` · actor [claude-code-f2616b4f](../actor/claude-code-f2616b4f.md)

## result

PASS

## statement

QA at 990e846 with a scratch build: a fresh manual session with the picked/typed transcript ran session done --json and printed "decisions": ["answer-<id>-1" .. "answer-<id>-4"]; the human output listed '- Proposed decision answer-<id>-N' for each; a repeat done --json returned already:true with no decisions field, and a no-transcript done had no decisions field. The describe-golden script passes, so the golden holds the new SessionDoneResponse.decisions field.

## Links

### Out

- verifies [done-reports-proposed](../capability/rhiza/lineage/interview-capture/done-reports-proposed.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
