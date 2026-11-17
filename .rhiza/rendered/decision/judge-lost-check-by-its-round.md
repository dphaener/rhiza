# judge-lost-check-by-its-round

decision · `01M7MGM9C244HF6K2Q88HWH25E`

## reasoning

The cause is in internal/run/checks.go pollCheck: on a resume, internal/api run_lost.go gives a check session that died unseen a lost exit, AgentStatus reports it Done with outcome fail, and pollCheck journals that judgement, so checkVerdict fails the check ('check failed: lost: ...') without reading the findings round and nextCheckActions opens cycle 2 over every check. Decision judge-sessions-ended-unseen already lets a lost agent step pass when the evidence of its work holds now (its until); a check's findings round, and a critic's verdicts, are that evidence for a checks step, recorded in the run's journal by rhiza run findings --check and rhiza run verdict --check. A lost exit's code is unknown, so without the recorded work it stays a fail, and known exits keep their meaning.

## statement

A check or critic session of a checks step that the run finds lost (an exit of state lost: it ended while no run watched it, as in a crash) is judged by the work it recorded rather than by its unknown exit: a lost check session is judged pass when the journal holds its findings round for that cycle, and a lost critic session is judged pass when the journal holds its verdict on every blocker and major finding of that round; checkVerdict then reads the round and the verdicts as for any passed session. A lost session without that recorded work, a session that exited with a known non-zero code, and a session that never read its start command are still judged fail. The done mark is not required, as a crash can come between rhiza run findings and rhiza session done.

## Links

### Out

- derived-from [plan-lost-check-sessions](../source/plan-lost-check-sessions.md)

### In

- decided-by [lost-check-judged-by-round](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step/lost-check-judged-by-round.md)

## History

- `01M7MGF8J9MA30XXQ8GTR4KXRW` Fix Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes
