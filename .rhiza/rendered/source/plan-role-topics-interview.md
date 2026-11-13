# plan-role-topics-interview

source · `01M79BEQC2C08ZJVASXVB26YVK`

## quotes

```json
[]
```

## raw_sha256

b881fb6685f6c5ef815b7664d2a79b5f0d7b4675bea670df420e523ded8c8da0

## session

3c4cbc82-631f-48a7-a942-01a1d2b2e5e4

## summary

Autonomous plan run for change set 01M76XRP (Plan role topics), serving m35-afk-prompts. No owner present. The planner read the change set intent, prompt-gap.md sections 2.1, 2.2, 4 (rec 15) and 5, the planned sibling Build role topics (01M76XR5, Decision build-role-topics-shape, its shared-block mechanism), the Leaf contract change set's intent (01M76XT3, which also names plan-review kinds), internal/guide/text/plan.md and plan-review.md, internal/guide/plan_test.go and review_test.go, internal/run/journal.go (Kinds, KindSeverity), internal/api/run_review.go, internal/api/guide_install.go and the builtin plan.toml and plan-milestone.toml prompts. It raised one question bundling four choices (delivery and the fate of the plan and plan-review topics, which change set adds the uncovered-intent and oversized-leaf kinds and at what severity, how the research pass and self-check work, how shared blocks and the severity block are reused). The arbiter answered with Decision plan-role-topics-shape: four new topics planner, plan-reviewer, plan-critic and plan-reviser installed as rhiza-planner, rhiza-plan-reviewer, rhiza-plan-critic and rhiza-plan-reviser; plan and plan-review stay as short indexes (not breaking); plan.toml and plan-milestone.toml prompts shrink to one line per role; this change set adds both kinds to run.Kinds as major, Leaf contract keeps no-failure-case and the review-check move; research pass and self-check are planner prose only; this change set depends on Build role topics and its shared blocks, and the severity block moves to internal/guide/text/shared/severities.md. The planner then wrote eight leaves under plan-workflow and linked the change set depends-on 01M76XR5, since that change set's nodes are not landed and a leaf cannot link to them.

## title

Plan role topics plan (autonomous)

## Links

### In

- derived-from [plan-role-topics-shape](../decision/plan-role-topics-shape.md)
- derived-from [plan-finding-kinds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-finding-kinds.md)
- derived-from [shared-severities-block](../capability/rhiza/workflows/built-in-workflows/plan-workflow/shared-severities-block.md)
- derived-from [planner-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-topic.md)
- derived-from [plan-reviewer-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-topic.md)
- derived-from [plan-critic-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-critic-topic.md)
- derived-from [plan-reviser-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-topic.md)
- derived-from [plan-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-role-skills-installed.md)
- derived-from [plan-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-prompts-one-line.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
