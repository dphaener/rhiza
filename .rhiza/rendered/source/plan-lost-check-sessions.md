# plan-lost-check-sessions

source · `01M7MGKJ4XGZR2S669Q1S2KMGA`

## quotes

```json
[]
```

## raw_sha256

30b716fbac86a2b9e041c83495822b59a0044dc96f2ebc3c7f0f82f382b3372c

## session

a5fc7c67-293d-49ba-a97a-01a1e907de70

## summary

Autonomous plan run for the bug 'Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes'. No owner was present and no question was raised: the code settles the cause and an earlier Decision (judge-sessions-ended-unseen) sets the pattern. Cause: on resume a check session that died unseen gets a 'lost' exit (internal/api/run_lost.go); AgentStatus reports it Done with outcome fail; pollCheck (internal/run/checks.go) journals that as the session's judgement, so checkVerdict fails the check with 'check failed: lost: ...' without reading its findings round, and nextCheckActions opens cycle 2 over every check. The agent-step precedent (engine.go, judge-sessions-ended-unseen) lets a lost session still pass when the evidence of its work holds now (the until of a done-waiting step). Plan: one leaf in internal/run/checks.go judging a lost check session by its recorded findings round for the cycle, and a lost critic by its recorded verdicts; a lost session without them, and every non-lost failed exit, still fail. The change set serves m4-runs-hold-up (a review run finishes without hand repair). Not breaking, not severe.

## title

Planning: a resumed classic review judges lost check sessions by their recorded work

## Links

### In

- derived-from [judge-lost-check-by-its-round](../decision/judge-lost-check-by-its-round.md)

## History

- `01M7MGF8J9MA30XXQ8GTR4KXRW` Fix Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes
