# plan-reviser-topic-verified

evidence · `01M79MR3Y0XCV6BFZDSSETQEXX` · actor [claude-code-8ff22888](../actor/claude-code-8ff22888.md)

## result

PASS

## statement

At 34af99081 re-ran the Checks block: go test -trimpath ./internal/guide/... (exit 0, includes TestPlanReviserText and TestListOrder), go test -trimpath ./cmd/rhiza -run TestScripts/guide (exit 0, guide ls counts 22 texts and lists plan-reviser), go test -trimpath ./internal/api -run Guide (exit 0), rhiza ship generate (exit 0, 0 changed), rhiza render (exit 0, no tracked diff). Read the diff: internal/guide/text/plan-reviser.md has the six sections in order; Must do names rhiza apply <file> --change <cs> --json and --record --json, splitting an oversized leaf into one-concern leaves, adding a leaf for an uncovered intent clause, never inventing an owner's words (missing quote becomes a question); Read first and Must not do carry the confirmed-major-only rule; command-budget, never-list and session-done blocks expand (checked in rhiza guide plan-reviser output). The test asserts all of these plus a negative case banning rhiza run findings/verdict, heredoc and findings/verdicts JSON. Nothing outside the leaf changed apart from the builder's actor record and its rendered file; no stubs or TODOs.

## Links

### Out

- verifies [plan-reviser-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-topic.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
