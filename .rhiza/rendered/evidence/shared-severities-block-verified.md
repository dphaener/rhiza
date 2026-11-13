# shared-severities-block-verified

evidence · `01M79MHTPGCD5D2D2R5KD16FKJ` · actor [claude-code-22c3fb0b](../actor/claude-code-22c3fb0b.md)

## result

PASS

## statement

At 228726468: go test -trimpath -count=1 ./internal/guide/... exit 0 (TestSeveritiesBlockShared passes; TestReviewSeverityBlockShared removed); rhiza ship generate exit 0; rhiza render exit 0 with no tree changes. Diff: severities block lives only in internal/guide/text/shared/severities.md, review.md and plan-review.md carry `<!-- shared: severities -->`; the test checks no unexpanded text holds ## Severities, review and plan-review show the expanded block, every topic naming the block shows it, uncovered-intent and oversized-leaf are major, and Expand of `<!-- shared: severity -->` errors naming severity. `rhiza guide review` shows the expanded section. No TODOs, no changes outside the leaf.

## Links

### Out

- verifies [shared-severities-block](../capability/rhiza/workflows/built-in-workflows/plan-workflow/shared-severities-block.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
