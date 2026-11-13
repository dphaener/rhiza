# plan-role-skills-installed-verified

evidence · `01M79NPJ0S33KGF8D375DP8Y1V` · actor [claude-code-5c24b65c](../actor/claude-code-5c24b65c.md)

## result

PASS

## statement

At 071bd5aab re-ran every Checks: command, each exit 0: go test -trimpath ./internal/guide/...; go test -trimpath ./internal/api -run GuideInstall (TestGuideInstallRoleSkills checks rhiza-planner, rhiza-plan-reviewer, rhiza-plan-critic and rhiza-plan-reviser written with their topic's description and expanded body, and unchanged on a second install); go test -trimpath ./internal/workflow/...; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install|describe-golden)$'; rhiza ship generate (0 changed); rhiza render (only this verifier's own actor record). The acceptance's go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' also passes: guide.txtar's guide ls lists planner, plan-reviewer, plan-critic and plan-reviser, and guide-install.txtar checks the four new skill paths in the install output. Read the diff: plan.md and plan-review.md are now short indexes naming each role topic as /rhiza-<role> and rhiza guide <role>, with rhiza session done; TestPlanIndexes checks both routes and keeps each index under TopicCap/3 (plan 204 tokens, plan-review 164). Plan's description in guide.go is unchanged and TestPlanKinds still checks its trigger phrases. Removed plan-review phrase tests are covered by the plan-reviewer, plan-critic and plan-reviser topic tests. Nothing changed outside the leaf, and no TODOs or stubs were left.

## Links

### Out

- verifies [plan-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-role-skills-installed.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
