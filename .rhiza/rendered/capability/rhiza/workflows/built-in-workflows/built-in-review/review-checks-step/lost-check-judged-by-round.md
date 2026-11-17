# lost-check-judged-by-round

capability · `01M7MGM9EZ9H0G8F45S2S629HK` · parent [review-checks-step](../review-checks-step.md)

## acceptance

New engine tests in internal/run/checks_test.go, run with go test -trimpath ./internal/run -run 'TestChecksStep|TestCritic|TestNextCheckActions': (1) a resumed checks step whose journal holds cycle 1 of three checks, each with its session and an empty findings round, and whose fake runner reports every session Done, Lost, outcome fail, passes in cycle 1: each check's judgement is pass, no cycle 2 entry is journalled and no check session is started; (2) a resumed step whose lost check session's round holds a major finding starts that check's critic, and a lost critic that recorded a verdict on that finding is judged pass and its verdict decides the check; (3) failure case: a lost check session with no findings round for the cycle is judged fail, and cycle 2 runs only that check; (4) failure case: a lost critic missing a verdict on a serious finding is judged fail and its check fails ('critic failed'); (5) failure case: a check session exiting with a known non-zero code (not lost) with a recorded round is still judged fail. The existing TestChecksStep* and TestCritic* tests pass unchanged.

## files

internal/run/checks.go, internal/run/checks_test.go

## must_not

Judge a lost session pass without its recorded round (check) or verdicts (critic); change how a non-lost exit, a never-submitted session, an agent step, a leaf fixer or verifier, or a change-set fixer is judged; change AgentStatus, Rejoin or the lost exit in internal/api; change the cycle cap or which checks a later cycle runs.

## statement

A checks step judges a lost check session by its recorded findings round for the cycle and a lost critic session by its recorded verdicts, so resuming a run whose check sessions died in a crash after recording their rounds passes or fails each check on its round instead of failing the cycle and re-running every check.

## Links

### Out

- decided-by [judge-lost-check-by-its-round](../../../../../../decision/judge-lost-check-by-its-round.md)

### In

- verifies [lost-check-judged-by-round-verified](../../../../../../evidence/lost-check-judged-by-round-verified.md)
- verifies [lost-check-judged-by-round-qa-1](../../../../../../evidence/lost-check-judged-by-round-qa-1.md)

## History

- `01M7MGF8J9MA30XXQ8GTR4KXRW` Fix Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes
