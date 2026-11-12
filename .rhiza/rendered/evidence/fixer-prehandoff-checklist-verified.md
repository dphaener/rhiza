# fixer-prehandoff-checklist-verified

evidence · `01M788DJK3F2QBTG3DBW2QBSMB` · actor [claude-code-dc5b5582](../actor/claude-code-dc5b5582.md)

## result

PASS

## statement

Verifier claude-code-dc5b5582 at 5f8bfd112: go test -trimpath ./internal/guide -run Review -v passes, incl. new TestReviewFixerPreHandoffChecklist, which asserts review.md's Fixer section names reproducing the findings, `rhiza ship generate`, `rhiza render`, the `Checks:` block and exit code; go test -trimpath ./internal/guide passes (review.md 799/800 tokens under the cap; to fit, the builder shortened unrelated review.md wording, e.g. the severity glosses and the merge-base git diff command). go test -trimpath ./internal/workflow passes, incl. new TestFixerPromptsNameChecksBlock, which asserts the review.toml fixer prompt and ship-change's ci-fix, gate-fix and conflict-fix prompts each name the `Checks:` block, `rhiza ship generate`, `rhiza render` and exit code. In the diff, gate-fix-merge and conflict-fix-landed get the same text. Caveat: make self-check in a clean detached checkout of 5f8bfd112 (RHIZA_ROOT/CHANGE/RUN unset) fails on a stale render: the claude-code-2be75af1 page is missing (its record came in with the run's evidence commit 4f44e4fb9), and this commit's index.md lists the untracked claude-code-14ab76b3. So the commit's 'make self-check: exit 0' held only in its dirty worktree. Commit a62c7cd14 on cs/pre-handoff-checklist already fixes this render, so no new fix is needed for this leaf.

## Links

### Out

- verifies [fixer-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
