# plan-reviewer-topic

capability · `01M79BH5P0K3C4P0VH374J64XE` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a plan-reviewer test that finds the six section headings in order, every kind in run.Kinds (so a kind added in Go without the topic fails), the vague-acceptance and unvalidated-assumption rules, `rhiza overlaps`, the line that `rhiza check` findings are not recorded, the findings heredoc and the empty-round rule, and the severities, session-done, command-budget and never-list blocks expanded; a negative case asserts it carries no `rhiza run verdict` or `rhiza apply` instructions.

## statement

`rhiza guide plan-reviewer` gives a fresh plan reviewer, on the six sections, afk's Crucible checklist in Rhiza's words: find the plan's change sets and read them through `rhiza slice` and `rhiza why`, record one finding per problem of the seven kinds (the five existing plus uncovered-intent and oversized-leaf, with a vague or unrunnable acceptance a missing-acceptance and a codebase assumption the code does not bear out a contradiction), run `rhiza overlaps` per change set, never record what `rhiza check` already reports, and record the round with the `rhiza run findings` heredoc, empty when clean; it carries the shared severities, session-done, command-budget and never-list blocks.

## Links

### Out

- depends-on [plan-finding-kinds](plan-finding-kinds.md)
- depends-on [shared-severities-block](shared-severities-block.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [plan-role-skills-installed](plan-role-skills-installed.md)
- verifies [plan-reviewer-topic-verified](../../../../../evidence/plan-reviewer-topic-verified.md)
- verifies [plan-reviewer-topic-qa-1](../../../../../evidence/plan-reviewer-topic-qa-1.md)
- verifies [plan-reviewer-topic-verified-2](../../../../../evidence/plan-reviewer-topic-verified-2.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
