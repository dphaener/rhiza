# fixer-prehandoff-checklist-verified-2

evidence · `01M7899X8NNS6N6B14YP3866JJ` · actor [claude-code-787640ac](../actor/claude-code-787640ac.md)

## result

PASS

## statement

At aa2f5da1c: go test -trimpath ./internal/guide -run Review passes, incl. TestReviewFixerPreHandoffChecklist, which slices review.md between '## Fixer' and '## Outside work' and asserts 'reproduce your findings', `rhiza ship generate`, `rhiza render`, '`Checks:` block' and 'exit code'; the whole guide package passes too (token caps). go test -trimpath ./internal/workflow passes, incl. TestFixerPromptsNameChecksBlock asserting the review.toml fixer prompt and ship-change.toml ci-fix, gate-fix and conflict-fix prompts each name the `Checks:` block, rhiza ship generate, rhiza render and exit code. Read the diff: review.md fixer section and the four prompts (plus gate-fix-merge and conflict-fix-landed) say to run reproducing commands, then rhiza ship generate and rhiza render before the one commit, with a Checks: block of commands and exit codes.

## Links

### Out

- verifies [fixer-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
