# plan-critic-topic-verified

evidence · `01M79N5K8AZVKYCPP3WE26G0YP` · actor [claude-code-85226b87](../actor/claude-code-85226b87.md)

## result

PASS

## statement

At 7a0a398f0 re-ran the Checks block: go test -trimpath ./internal/guide/... ok (TestPlanCriticText: six sections in order, rhiza run status --json in Read first, verdict heredoc with confirmed/downgraded+severity/dropped, never add/raise rules, command-budget/never-list/severities/session-done expanded, negative case bans rhiza run findings, rhiza apply and a findings heredoc), go test ./cmd/rhiza -run TestScripts/guide ok, go test ./internal/api -run Guide ok, rhiza ship generate exit 0 (0 changed), rhiza render exit 0 with a clean tree. Read the diff: only the plan-critic topic, its registration, list-order/guide.txtar counts and the build's log/rendered entries; no stubs or TODOs; rendered `rhiza guide plan-critic` matches the run verdict input schema.

## Links

### Out

- verifies [plan-critic-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-critic-topic.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
