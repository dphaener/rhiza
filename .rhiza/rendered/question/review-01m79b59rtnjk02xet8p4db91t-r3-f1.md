# review-01m79b59rtnjk02xet8p4db91t-r3-f1

question · `01M79CV2GRBTKSEB6Y95AYYZJ9`

## statement

Plan review round 3 left finding f1 (contradiction on plan-guide): plan-guide (updated here) says the planner topic carries the apply recipe and the no-plan-file/no-checkboxes rule, and its acceptance checks `rhiza guide planner` holds an apply recipe for a leaf and says no plan file and no checkboxes are kept. But Decision plan-role-topics-shape and leaf planner-topic move the apply recipes, no plan file and who records a pick into the separate `plan-recipes` reference topic (planner-topic's acceptance asserts those phrases live in plan-recipes). A builder cannot satisfy both: either planner-topic's acceptance or plan-guide's fails. plan-guide's statement and acceptance should name plan-recipes for the recipe and no-plan-file rule. The critic confirmed it: plan-guide's acceptance requires `rhiza guide planner` to hold an apply recipe for a leaf (statement, acceptance, decided-by) and the no-plan-file/no-checkboxes rule, while planner-topic's statement and acceptance move the apply recipe and no plan file into plan-recipes and keep planner within TopicCap; Decision plan-role-topics-shape puts the recipes in plan-recipes. A builder must either break the split or fail plan-guide's acceptance, so the leaves contradict each other. How should it be settled?

## Links

### Out

- blocks [plan-guide](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide.md)

### In

- answers [plan-guide-follows-recipes-split](../decision/plan-guide-follows-recipes-split.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
