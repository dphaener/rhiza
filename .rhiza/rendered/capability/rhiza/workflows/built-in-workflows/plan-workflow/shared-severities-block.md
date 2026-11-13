# shared-severities-block

capability · `01M79BH58V0BQMAAJ8RS4AJ9MT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with TestReviewSeverityBlockShared replaced by a test that the unexpanded texts hold the `## Severities` section only in shared/severities.md, every topic naming the block shows it expanded (review at least), the block names uncovered-intent and oversized-leaf as major, and (failure case) a text naming a misspelt block such as `<!-- shared: severity -->` fails expansion with an error naming it.

## statement

The severities block (blocker, major, minor, nit, the per-check mapping, and the plan-review mapping with missing-acceptance, contradiction, outside-intent, uncovered-intent and oversized-leaf major and the rest minor) lives once in internal/guide/text/shared/severities.md, and each topic that carries it (review, and the plan reviewer and critic topics) names it with a `<!-- shared: severities -->` line instead of its own copy.

## Links

### Out

- depends-on [plan-finding-kinds](plan-finding-kinds.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [plan-reviewer-topic](plan-reviewer-topic.md)
- depends-on [plan-critic-topic](plan-critic-topic.md)
- verifies [shared-severities-block-verified](../../../../../evidence/shared-severities-block-verified.md)
- verifies [shared-severities-block-qa-1](../../../../../evidence/shared-severities-block-qa-1.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
