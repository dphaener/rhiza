# review-role-skills-installed-verified-5

evidence · `01M79WAASK4B5GG3GS11QK7VS9` · actor [claude-code-dcf94a53](../actor/claude-code-dcf94a53.md)

## result

PASS

## statement

At 7d6a3b395 (attempt 5): re-ran the Checks block: go build ./cmd/rhiza exit 0; rhiza-x guide check-compliance/check-quality/check-model each has 1 'fakes' line; go test -trimpath ./internal/guide/... ok; rhiza ship generate 0; rhiza render 0 with a clean tree. Acceptance: go test -trimpath ./internal/api -run GuideInstall ok; go test -trimpath ./internal/guide/... ok (review_test asserts the nine topics with their routes, scope routed to check-compliance, the index under TopicCap/3); go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' PASS. Diff only adds the fakes-only rule to three check topics with tests and trims two untested lines; no stubs or TODOs.

## Links

### Out

- verifies [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
