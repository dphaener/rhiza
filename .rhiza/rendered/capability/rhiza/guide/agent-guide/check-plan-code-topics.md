# check-plan-code-topics

capability · `01M7AHW2WHSC1M6D5AFSGEJRFR` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run '^(TestCheckPlanText|TestCheckPlanNamesNoRepoCommand|TestCheckCodeText|TestCheckCodeNamesNoRepoCommand|TestReviewIsAnIndex|TestBuildReviewCommandBudget)$' shows check-plan holding the compliance and model rules, check-code holding the quality rules, both with the role sections and shared blocks and under the size cap, and the review topic naming check-plan, check-code and check-qa for a conductor run and the old four for the classic workflows; go test -trimpath ./internal/api -run TestGuideInstallRoleSkills shows skills/rhiza-check-plan/SKILL.md and skills/rhiza-check-code/SKILL.md installed beside the old role skills. Touches internal/guide/text/check-plan.md, internal/guide/text/check-code.md, internal/guide/text/review.md, internal/guide/guide.go, internal/guide/check_plan_test.go, internal/guide/check_code_test.go, internal/guide/review_test.go, internal/guide/budget_test.go, internal/api/guide_install.go, internal/api/guide_install_test.go.

## statement

Rhiza ships two review role topics and skills: check-plan (the compliance and model checks' rules merged: each leaf's acceptance met as written, under-built, scope creep, unaccounted files, and code that contradicts the slice) and check-code (the quality check's rules under the new name); check-compliance, check-model and check-quality stay unchanged for the classic ship-change, review and docs-review, and the review index topic names which runs load which topics.

## Links

### Out

- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)
- decided-by [conductor-three-checks](../../../../decision/conductor-three-checks.md)

### In

- depends-on [conduct-skill](conduct-skill.md)
- verifies [check-plan-code-topics-verified](../../../../evidence/check-plan-code-topics-verified.md)
- verifies [check-plan-code-topics-qa-1](../../../../evidence/check-plan-code-topics-qa-1.md)
- verifies [check-plan-code-topics-qa-2](../../../../evidence/check-plan-code-topics-qa-2.md)
- verifies [pr-575-a-rhiza-conduct-skill-conducts](../../../../evidence/pr-575-a-rhiza-conduct-skill-conducts.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
