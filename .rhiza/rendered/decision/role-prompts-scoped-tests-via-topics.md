# role-prompts-scoped-tests-via-topics

decision · `01M79BP1HNWZ98HJS9C16QZ95S`

## reasoning

Read the question, role-prompts-scoped-tests (statement and acceptance requiring each prompt to carry the rule), leaf review-prompts-one-line (one-line prompts with no test-scope text), decision review-role-topics-shape (prompts shrink to one line naming their topic; everything else moves into the topics through shared blocks, the test-scope block included in every check, critic and fixer topic per their leaves' acceptance), and internal/workflow/workflow_test.go:1552 TestBuiltinRolePromptsScopeTests. The change set's intent is to move role instructions into topics, and the check-*/critic/fixer leaves already require the shared test-scope block, so the rule still reaches every agent; only the older node's wording (where the rule lives) is stale. Re-adding the rule to the prompts would contradict the shape decision and the one-line leaf. Phrasing the node as 'in the prompt or through the named topic' keeps the build workflow's prompts covered whichever form they take (the sibling Build role topics decision also shrinks them), and keeps the substance (scoped tests, no named commands) testable. Not breaking: no command, flag, --json shape, setting or .rhiza/ format changes.

## statement

Settle f1 by updating role-prompts-scoped-tests in this change set, not by putting test-scope text back into the one-line prompts: its statement becomes that every builder, verifier, check, critic and fixer of the built-in build and review workflows is told to run only the tests covering what its change touches and never the repo's full suite (which the PR's CI runs), either in its prompt or through the guide topic its prompt names, whose text includes the shared test-scope block; the review workflow's check, critic and fixer prompts reach it through check-<check>, critic and fixer. Its acceptance changes accordingly: the internal/workflow test (TestBuiltinRolePromptsScopeTests) asserts each such prompt either carries the rule or names a `rhiza guide <topic>` whose topic text (with shared blocks expanded) carries it, and that no prompt names go test, make check, make test or testscript; go test -trimpath ./internal/workflow/... and ./internal/guide/... pass. Leaf review-prompts-one-line stays as planned; the builder fixes the test with the node.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79a41352xbz4fc2fjt5903j-r3-f1](../question/review-01m79a41352xbz4fc2fjt5903j-r3-f1.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
