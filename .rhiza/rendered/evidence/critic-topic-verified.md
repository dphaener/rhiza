# critic-topic-verified

evidence · `01M79NHXW9CMRCH4R6NA0GWHZ8` · actor [claude-code-c290ab42](../actor/claude-code-c290ab42.md)

## result

PASS

## statement

At c8a122646 re-ran the Checks: block: go test -trimpath ./internal/guide/... (exit 0), go test -trimpath ./internal/api -run Guide (exit 0), go test -trimpath ./cmd/rhiza -run TestScripts/guide (exit 0), rhiza ship generate (exit 0, 0 changed), rhiza render (exit 0). Read the diff: internal/guide/text/critic.md adds the topic on the six sections in order with the session-done, command-budget, never-list and severity shared blocks expanded (checked via a scratch build's rhiza guide critic), the isolation rule in Task with the own-check-last-round exception, confirmed/downgraded (with a lower severity)/dropped, the never add/raise/rewrite rule and the rhiza run verdict --check heredoc (matches rhiza describe run.verdict). internal/guide/critic_test.go TestCriticText asserts each acceptance point; TestCriticRecordsNoFindingAndNoCommit is the negative case (no rhiza run findings, no git commit/Checks:/trailers). Other changes are the registry entry, list-order and guide ls counts, a severity-test scope tweak for the critic topic, and the builder's actor record; nothing outside the leaf, no TODOs or stubs.

## Links

### Out

- verifies [critic-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/critic-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
