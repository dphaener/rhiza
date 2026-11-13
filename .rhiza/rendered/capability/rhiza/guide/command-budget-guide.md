# command-budget-guide

capability · `01M76XY4EE4QA5EKNQXHR7R9XH` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes: TestBuilderText, TestVerifierText and TestSharedPhrasesInExactlyOneFile assert the shared command-budget block holds the budget rule (no sleep or poll, about four minutes per command, background once and wait once) and that the builder and verifier topics carry it inside Must do; TestBuildReviewCommandBudget loops over the nine expanded review role topics and asserts each holds the rule and that its Command budget section names no Go, make or Rhiza-repo command.

## statement

The builder and verifier guide topics and every review role topic (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix) tell builders, verifiers, checks, critics and fixers never to sleep or poll, never to run a single command longer than about four minutes, and to run a long build or test once in the background and wait for it once, in language-neutral words; each carries it as the shared command-budget block (the builder and verifier inside Must do), and the review index does not hold it.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../decision/efficient-cycle-milestone-shape.md)

### In

- blocks [review-01m79a41352xbz4fc2fjt5903j-r3-f3](../../../question/review-01m79a41352xbz4fc2fjt5903j-r3-f3.md)
- verifies [command-budget-guide-verified](../../../evidence/command-budget-guide-verified.md)
- verifies [command-budget-guide-qa-1](../../../evidence/command-budget-guide-qa-1.md)
- verifies [pr-501-command-budget](../../../evidence/pr-501-command-budget.md)
- verifies [pr-521-build-role-topics](../../../evidence/pr-521-build-role-topics.md)
- verifies [pr-537-review-role-topics](../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XQX53GGNWVCQHXVV510NH` Builder, verifier, check and fixer guidance says: never sleep or poll, no single command longer than about four minutes, and a long command runs once in the background and is waited for once (afk execute.js rule). Evidence: verifier session p90 33.7m against a 3.7m median, builder p90 25-36m against about 4m (prompt-gap rec 6).
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
