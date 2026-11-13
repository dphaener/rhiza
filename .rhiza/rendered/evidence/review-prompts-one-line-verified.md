# review-prompts-one-line-verified

evidence · `01M79Q6GEC40V0GW0D6R2N3QYV` · actor [claude-code-f34688f1](../actor/claude-code-f34688f1.md)

## result

PASS

## statement

At bba82e799 re-ran the Checks: block and acceptance: go test -trimpath ./internal/workflow/... ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(review-workflow|review-checks-step|autonomous-review|review-summary|describe-golden)' ok; rhiza ship generate and rhiza render exit 0 with no tree changes; make self-check exit 0 (twice in a clean detached worktree; one earlier run in the run checkout exited 2 and passed on immediate rerun with no change). Diff read: builtin review.toml has checks compliance, quality, qa, model with the exact one-line /rhiza-check-<check> prompts, critic and fixer one-line prompts; docs-review.toml drops scope, shares check and critic lines byte for byte, fixer is the builtin line plus the humanize/doc-prose-humanized/doc-examples-plain-commands clause. TestReviewBuiltin asserts exact prompts and absence of session done, Checks:, full suite; new TestDocsReviewMatchesBuiltinReview reads docs-review.toml and asserts no scope, equal check/critic prompts, fixer prefix and doc clause names. review-workflow.txtar fake updated for 4 checks; other scripts had no scope expectations; describe-golden unchanged. Other files are rhiza render output and the builder's actor record. No stubs or TODOs.

## Links

### Out

- verifies [review-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-prompts-one-line.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
