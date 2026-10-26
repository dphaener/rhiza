# interview-capture-scenario

evidence · `01M5WE456ARA8S7HZ7E2NAX6ZZ`

## result

pass

## statement

The session-interview and session-run scenarios propose one Decision per answered question and skip rejected ones and repeats.

## test

cmd/rhiza/testdata/script/session-interview.txtar

## Links

### Out

- verifies [interview-capture](../capability/rhiza/lineage/interview-capture.md)

## History

- `01M5SZJ2B3PSDRVX736ZQX0CAZ` Interview answers as Decisions: after an agent step, Genesis reads the session and proposes a Decision for each answer, a picked option as the agent's wording accepted by the owner and only typed text quoted as the owner's.
