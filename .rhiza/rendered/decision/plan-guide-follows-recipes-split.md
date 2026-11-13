# plan-guide-follows-recipes-split

decision · `01M79CWAVQCJY9J3Q69R22D0Y9`

## reasoning

Read the question, plan-guide, planner-topic and Decision plan-role-topics-shape. The Decision, which decides both leaves, puts the apply recipes and who records a pick in the plan-recipes reference topic because the planner part of today's plan topic (about 2000 tokens) cannot fit one topic under TopicCap 1000; planner-topic's statement and acceptance carry that split, asserting the apply recipe and no plan file move to plan-recipes. plan-guide was only updated to point at the new topics and still names the planner topic for the recipe, which is the older wording. Moving the recipe back into planner would break TopicCap and the Decision; rewording plan-guide costs nothing and keeps the content (the rule is still in the guide, one topic away, named from planner's Read first). The review finding itself proposes this fix and the critic confirmed it. No CLI, flag or file format changes, so no breaking note.

## statement

Settle the contradiction in favour of the split: plan-guide follows Decision plan-role-topics-shape and leaf planner-topic. The apply recipe for a leaf (statement and acceptance fields, a decided-by link), the no-plan-file and no-checkboxes rule, and who records a pick live in `rhiza guide plan-recipes`, not in `rhiza guide planner`. plan-guide's statement and acceptance are reworded to say so: the planner topic carries what a planning session plans, in what order, and that every leaf, acceptance check, Decision and question goes into the model through rhiza apply, and its Read first names plan-recipes and plan-interview; plan-guide's acceptance checks that rhiza guide plan names the four role topics, that rhiza guide planner names plan-recipes, and that rhiza guide plan-recipes holds the leaf apply recipe and says no plan file and no checkboxes are kept. planner-topic and the Decision are unchanged.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79b59rtnjk02xet8p4db91t-r3-f1](../question/review-01m79b59rtnjk02xet8p4db91t-r3-f1.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
