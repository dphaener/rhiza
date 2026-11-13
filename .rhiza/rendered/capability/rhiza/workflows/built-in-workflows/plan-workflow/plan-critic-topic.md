# plan-critic-topic

capability · `01M79BH5XRV6CSB9GDW3VGA5JT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a plan-critic test that finds the six section headings in order, `rhiza run status --json`, the verdict heredoc with confirmed, downgraded (with its severity) and dropped, the never-add and never-raise rule, and the severities, session-done, command-budget and never-list blocks expanded; a negative case asserts it carries no `rhiza run findings` heredoc or `rhiza apply` instructions.

## statement

`rhiza guide plan-critic` gives a plan critic, on the six sections: read the latest round in `rhiza run status --json`, check each finding against the model and the code, and confirm it, downgrade a blocker or major one to a lower severity, or drop it, recorded with the `rhiza run verdict` heredoc; never add a finding, raise a severity or rewrite one; with the shared severities, session-done, command-budget and never-list blocks.

## Links

### Out

- depends-on [shared-severities-block](shared-severities-block.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [plan-role-skills-installed](plan-role-skills-installed.md)
- verifies [plan-critic-topic-verified](../../../../../evidence/plan-critic-topic-verified.md)
- verifies [plan-critic-topic-qa-1](../../../../../evidence/plan-critic-topic-qa-1.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
