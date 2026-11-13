# review-role-skills-installed-verified-8

evidence · `01M79Y0TF5PR7QQCVNF6X83CBN` · actor [claude-code-e5d99856](../actor/claude-code-e5d99856.md)

## result

PASS

## statement

Attempt 8 at 895129820: go test -trimpath ./internal/api -run GuideInstall ok; go test -trimpath ./internal/guide/... ok; go test -trimpath ./internal/workflow/... ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' ok. A scratch build's rhiza guide ls lists check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix; rhiza guide review is a 253-token index naming all nine with their /rhiza-* skill and rhiza guide routes and sends a check named scope to check-compliance; rhiza guide critic shows the role-neutral 'Run only the tests covering what the change touches' sentence that decision test-scope-block-names-every-role asked for (shared/test-scope.md, shared_test.go, review_test.go updated). make self-check exits 0 in a clean detached worktree with the run env unset (in the run env it fails only on this verifier's own unrendered actor record). rhiza ship generate and rhiza render exit 0 and leave the tree clean. No TODOs or stubs in the diff; no changes outside the leaf.

## Links

### Out

- verifies [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
