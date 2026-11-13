# plan-role-skills-installed

capability · `01M79BH6B4G10X94QE56HC9RS4` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/api -run GuideInstall passes with the four new skill files written with their topic's description and expanded body, and a second install reporting them unchanged; go test -trimpath ./internal/guide/... asserts the plan and plan-review indexes name the four topics with both routes, plan's description is unchanged, and each index stays under a third of the TopicCap; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' passes with the new topics in guide ls and the new skill paths in the install output.

## statement

`rhiza guide ls` lists the planner, plan-reviewer, plan-critic and plan-reviser topics, guide install writes them as the rhiza-planner, rhiza-plan-reviewer, rhiza-plan-critic and rhiza-plan-reviser skills beside the others, and the plan topic (still installed as rhiza-plan, its description unchanged so the skill still triggers) and the plan-review topic are each a short index naming the four role topics, their skill and `rhiza guide` routes and which role loads which, so workflows and chats that still say /rhiza-plan or `rhiza guide plan-review` keep working.

## Links

### Out

- depends-on [planner-topic](planner-topic.md)
- depends-on [plan-reviewer-topic](plan-reviewer-topic.md)
- depends-on [plan-critic-topic](plan-critic-topic.md)
- depends-on [plan-reviser-topic](plan-reviser-topic.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [plan-prompts-one-line](plan-prompts-one-line.md)
- depends-on [plan-guide](plan-guide.md) — Its acceptance needs the plan index and the four role topics (planner's apply recipe included), which plan-role-skills-installed and the topics it depends on build; built after them, this leaf only checks and tightens that work.
- depends-on [plan-review-guide](plan-review-guide.md) — Its acceptance needs the plan-review index and the plan-reviewer, plan-critic and plan-reviser topics, which plan-role-skills-installed and the topics it depends on build; built after them, this leaf only checks and tightens that work.
- verifies [plan-role-skills-installed-verified](../../../../../evidence/plan-role-skills-installed-verified.md)
- verifies [plan-role-skills-installed-qa-1](../../../../../evidence/plan-role-skills-installed-qa-1.md)
- verifies [plan-role-skills-installed-verified-2](../../../../../evidence/plan-role-skills-installed-verified-2.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
