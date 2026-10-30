# done-reports-proposed-verified

evidence · `01M658SAG12HS4PJW1WNHK85Y9` · actor [claude-code-1944333a](../actor/claude-code-1944333a.md)

## result

PASS

## statement

At 5e03712: go test ./cmd/rhiza -run 'TestScripts/(session-done-captures|describe-golden|session-run-tokens)$' -count=1 passes; session-done-captures.txtar checks that session done --json prints "decisions": [ with the answer-158fe87f-N handles, and that the field is absent on a repeat done. Running describe-golden with -update leaves no diff, and the golden holds SessionDoneResponse.decisions ("Handles of the Decisions proposed from the questions the owner answered"). TestSessionDoneReportsProposed (internal/api) passes and checks both handles in the response and the human output. session-run-tokens (c) now expects the capture warning before the token warning, with no decisions field, which fixes attempt 1's finding. make check exits 0 and leaves the tree clean.

## Links

### Out

- verifies [done-reports-proposed](../capability/rhiza/lineage/interview-capture/done-reports-proposed.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
