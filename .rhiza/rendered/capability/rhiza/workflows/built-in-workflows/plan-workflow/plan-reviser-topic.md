# plan-reviser-topic

capability · `01M79BH640TST58REPF7FSDXJ5` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a plan-reviser test that finds the six section headings in order, `rhiza apply` with `--change` and `--record`, the split-an-oversized-leaf and leaf-for-an-uncovered-clause rules, the confirmed-major-only rule, the never-invent-a-quote rule, and the session-done, command-budget and never-list blocks expanded; a negative case asserts it carries no findings or verdict heredoc.

## statement

`rhiza guide plan-reviser` gives a plan reviser, on the six sections: fix each confirmed major finding of the latest round in the model and nothing else, one `rhiza apply <file> --change <cs> --json` per change set (`--record` for a node in a landed record), splitting an oversized leaf into leaves of one concern each and adding a leaf for an uncovered intent clause; never invent an owner's words (a missing quote becomes a question); with the shared session-done, command-budget and never-list blocks.

## Links

### Out

- depends-on [plan-finding-kinds](plan-finding-kinds.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [plan-role-skills-installed](plan-role-skills-installed.md)
- verifies [plan-reviser-topic-verified](../../../../../evidence/plan-reviser-topic-verified.md)
- verifies [plan-reviser-topic-qa-1](../../../../../evidence/plan-reviser-topic-qa-1.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
