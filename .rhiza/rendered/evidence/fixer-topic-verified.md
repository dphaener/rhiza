# fixer-topic-verified

evidence · `01M79NT7G1HQ192AT460DSHYQX` · actor [claude-code-0a7ab413](../actor/claude-code-0a7ab413.md)

## result

PASS

## statement

At 436b8f29 re-ran the Checks: block: go test -trimpath ./internal/guide/... (TestFixerText, TestFixerNamesNoRepoCommand, TestListOrder pass with -count=1), go test -trimpath ./internal/api -run Guide, go test -trimpath ./cmd/rhiza -run TestScripts/guide, rhiza ship generate, rhiza render: all exit 0, tree clean. Diff read: fixer.md has the six sections in order, done/test-scope/unstable-ci/command-budget/outside-work/never-list/severity shared blocks expanded (no markers left in a built `rhiza guide fixer`), fix-only-named, reproduce-first, rhiza ship generate, rhiza render, Checks: above the single final trailer paragraph, Leaf:/Review: forms, unfixed-never-dropped and never push/rebase/force/--no-verify, each asserted by TestFixerText; TestFixerNamesNoRepoCommand is the negative case for Go/make/repo commands. Only guide files, guide.txtar count and the builder's actor record changed; no stubs or TODOs.

## Links

### Out

- verifies [fixer-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
