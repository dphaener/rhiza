# review-01m79a41352xbz4fc2fjt5903j-r3-f1

question · `01M79BMA02G94Q32AF1ZHX5TEK`

## statement

Plan review round 3 left finding f1 (contradiction on role-prompts-scoped-tests): Says every check, critic and fixer prompt of the builtin review workflow carries the scoped-tests rule (and its internal/workflow test asserts it), but leaf review-prompts-one-line makes those prompts one line with no test-scope text. The plan updates neither this node nor its acceptance, so the built change contradicts the model and fails that test; update it to say the prompts reach the rule through their topics (the shared test-scope block). The critic confirmed it: role-prompts-scoped-tests says every check, critic and fixer prompt of the builtin review workflow carries the scoped-tests rule, and internal/workflow/workflow_test.go:1556 asserts it; leaf review-prompts-one-line's acceptance requires check prompts with no test-scope text and critic/fixer prompts equal to one-line forms without it. No plan change set updates the node, so the landed acceptance fails. Major stands. How should it be settled?

## Links

### Out

- blocks [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)

### In

- answers [role-prompts-scoped-tests-via-topics](../decision/role-prompts-scoped-tests-via-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
