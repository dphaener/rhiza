# plan-reviewer-topic-verified

evidence · `01M79N01Y45JK9V7XZCGZQ46XF` · actor [claude-code-74caff14](../actor/claude-code-74caff14.md)

## result

PASS

## statement

At 3bf8cdfaa re-ran the Checks: block: go test -trimpath ./internal/guide/... (exit 0, includes the acceptance command), go test -trimpath ./cmd/rhiza -run TestScripts/guide (exit 0), go test -trimpath ./internal/api -run Guide (exit 0), rhiza ship generate (exit 0), rhiza render (exit 0, no tracked changes). Read the diff: internal/guide/text/plan-reviewer.md has the six sections in order, all seven run.Kinds, the vague/unrunnable acceptance rule as missing-acceptance, the unvalidated codebase assumption as contradiction, rhiza overlaps per change set, the line that rhiza check findings are not recorded, the rhiza run findings heredoc and the empty-round rule, and expands the command-budget, never-list, severities and session-done blocks; TestPlanReviewerText asserts each of these (kinds iterated from run.Kinds) plus a negative case for rhiza run verdict, rhiza apply and verdicts. guide.go, guide_test.go and guide.txtar register the topic; only the leaf's files and its actor record changed; no stubs or TODOs.

## Links

### Out

- verifies [plan-reviewer-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-topic.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
