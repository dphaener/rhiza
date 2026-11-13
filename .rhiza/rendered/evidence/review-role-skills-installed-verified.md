# review-role-skills-installed-verified

evidence · `01M79PQZFXHESFSSPJSHBWPBVG` · actor [claude-code-199a84d1](../actor/claude-code-199a84d1.md)

## result

PASS

## statement

At 300a91dd4 re-ran the Checks: block and acceptance: go test -trimpath ./internal/guide/... (TestReviewIsAnIndex: index names all nine topics as `/rhiza-<t>` (`rhiza guide <t>`), routes scope to check-compliance, under TopicCap/3; TestReviewRoleText finds the old review.md phrases in the role topics; budget_test covers the nine topics) exit 0; go test -trimpath ./internal/api -run GuideInstall (TestGuideInstallReviewRoles: nine skill files equal frontmatter with topic description plus body, second install reports them unchanged) exit 0; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' (guide ls lists the nine topics, install output lists the nine new skill paths, 19 written) exit 0; rhiza ship generate exit 0 (0 changed); rhiza render exit 0. Diff read: guide_install.go writes the nine role skills, review.md is a short index, only leaf files and the builder's actor record changed, no TODOs or stubs; tests use fakes only.

## Links

### Out

- verifies [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
