# plan-workflow-autonomous-flake-plan

source · `01M7RFPWW4DYR119D54RH9S6T1`

## quotes

```json
[]
```

## raw_sha256

1723fb143f4bba294b0c22d1908faaa23bba50e177d5809920e9e4dddbdcc0fb

## session

7370cbac-6687-452c-a560-01a1f0f6f80c

## summary

Autonomous plan run, no owner. CI failure read from the failed first attempt of PR #655's CI run (scripts (1/4), head 1be14770b): plan-workflow-autonomous.txtar:68 found the journal order 'started process entered session blocked asked arbiter arbitrated answered unblocked tokens tokens ended tokens entered ...' and its regex allows tokens entries only before the interview step's ended entry. Cause: a run session's token use is journaled by three writers (internal/api/session_done.go at session done, stepRunner.AgentEnded in internal/api/run.go just before the engine writes ended, and internal/api/session_run.go when the agent process exits, a separate process). The exit-time write races the engine and may land after ended; recordTokenUse documents 'the later reading wins', and the journal fold (internal/run/journal.go, TypeTokens) is order-independent, so the product behaves as designed and the script's ordering assumption is wrong. Reproduced at main f7313222e with injected delays only (no stress): a temporary time.Sleep(1500ms) after the engine's ended entry (internal/run/engine.go) and a sleep 0.8 between the fake planner's rhiza session done and its exit; 3 of 3 runs failed at line 68 with the exact CI sequence; without the delays the order was tokens tokens tokens ended (3 passes). Both edits were reverted. Planned: fix the script's assertion (drop tokens from the order check, assert the interview session's tokens entry separately); do not quarantine; no product change.

## title

Plan of the plan-workflow-autonomous flake (autonomous research)

## Links

### In

- derived-from [plan-autonomous-tokens-order-free](../decision/plan-autonomous-tokens-order-free.md)

## History

- `01M7RFDAZ53YZ3775D7SV70Z4G` Fix flaky test TestScripts/plan-workflow-autonomous: failed then passed on a rerun in PR #655 at head 1be14770b82366ab1535a1b38a36a890e8b75939; reproduce it locally, repeatedly, then fix it or quarantine it
