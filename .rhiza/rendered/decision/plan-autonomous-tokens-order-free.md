# plan-autonomous-tokens-order-free

decision · `01M7RFQFV4PXDA50Y7RHTJPX32`

## reasoning

recordTokenUse is documented as 'the later reading wins' with three writers (session done, AgentEnded, agent exit); the exit-time write from internal/api/session_run.go races the engine's ended entry by design, and the TypeTokens fold in internal/run/journal.go keys by session, not position. Injected delays reproduced the CI order 3 of 3.

## statement

Fix the plan-workflow-autonomous flake in its script, not by quarantine and not in Go: a run session's tokens journal entries may land in any order relative to its step's ended entry (the agent's exit-time token reading is written by its own process, and the journal fold is order-independent), so the script's journal-order assertion stops placing tokens entries and checks them separately

## wording

Fix the script's journal-order assertion; do not quarantine; no product change. Planner's ruling in an autonomous plan run, from the reproduction

## Links

### Out

- derived-from [plan-workflow-autonomous-flake-plan](../source/plan-workflow-autonomous-flake-plan.md)

### In

- decided-by [plan-autonomous-journal-order-ignores-tokens](../capability/plan-autonomous-journal-order-ignores-tokens.md)

## History

- `01M7RFDAZ53YZ3775D7SV70Z4G` Fix flaky test TestScripts/plan-workflow-autonomous: failed then passed on a rerun in PR #655 at head 1be14770b82366ab1535a1b38a36a890e8b75939; reproduce it locally, repeatedly, then fix it or quarantine it
