# plan-review-guide

capability · `01M5WSGD6GSZAS9X5C75R728DG` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go guide tests: rhiza guide ls lists plan-review, plan-reviewer, plan-critic and plan-reviser; plan-review names the three role topics; plan-reviewer names every kind in run.Kinds and rhiza run findings with its format and how the plan's change sets are found; plan-critic names rhiza run verdict with its format and the rule that a critic never adds a finding; each topic stays within TopicCap.

## statement

rhiza guide plan-review is an index naming the plan-reviewer, plan-critic and plan-reviser topics: the reviewer, given only a change set ID, finds the plan's change sets (that one and every change set serving an intent it adds), reads them through rhiza slice and rhiza why and records findings of the plan-review kinds with rhiza run findings; the critic confirms, downgrades or drops each finding with rhiza run verdict, never adding one; the reviser fixes the confirmed findings in the model.

## Links

### Out

- depends-on [plan-review-journal](plan-review-journal.md)
- depends-on [plan-role-skills-installed](plan-role-skills-installed.md) — Its acceptance needs the plan-review index and the plan-reviewer, plan-critic and plan-reviser topics, which plan-role-skills-installed and the topics it depends on build; built after them, this leaf only checks and tightens that work.
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)

### In

- depends-on [plan-workflow-file](plan-workflow-file.md)
- verifies [plan-review-guide-verified](../../../../../evidence/plan-review-guide-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
