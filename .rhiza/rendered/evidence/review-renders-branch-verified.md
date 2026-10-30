# review-renders-branch-verified

evidence · `01M66DG22ATH1JX4NHZ8X79JW0` · actor [claude-code-863ddf1a](../actor/claude-code-863ddf1a.md)

## result

PASS

## statement

At e7f894a, review.toml adds a render step (review -> render -> signoff; request-changes returns through review and render) that runs rhiza render and commits .rhiza/rendered and .rhiza/log as 'Render after the review' with a Change: trailer, or commits nothing when there is nothing to render. Ran go test ./cmd/rhiza -run 'TestScripts/(review-workflow|build-workflow|run-ship-change)' -count=1 and go test ./internal/workflow -count=1: both ok. review-workflow.txtar runs a review whose qa check and alpha fixer/verifier write records. It asserts that the branch log ends with 'Render after the review' after 'chore(rhiza): evidence for alpha', that git status --porcelain -- .rhiza is empty, and that rhiza check passes on the branch. make check (fmt-check, vet, all tests, self-check) exit 0.

## Links

### Out

- verifies [review-renders-branch](../capability/rhiza/workflows/review-renders-branch.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
