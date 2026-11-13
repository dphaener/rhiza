# guides-scoped-tests

capability · `01M6M0EPMPD72KB5AZWBB9MRGM` · parent [guide](../guide.md)

## acceptance

TestBuilderText and TestVerifierText in internal/guide assert the builder and verifier topics carry the shared test-scope block, TestSharedPhrasesInExactlyOneFile that the block states the rule for builders and verifiers, and TestReviewText that review.md states it for checks, critics and fixers; TestBuilderText and TestReviewText assert the builder and review texts contain no go test, make or testscript; go test -trimpath ./internal/guide passes

## statement

The builder, verifier and review guide topics, and so the rhiza-builder, rhiza-verifier and rhiza-review skills, tell builders, verifiers, checks, critics and fixers to run only the tests covering what their change touches, with the repo's own scoped commands, and never the repo's full suite, because the PR's CI runs it before merge; the rule names no Go, make or testscript command

## Links

### Out

- decided-by [ci-runs-full-suite](../../../decision/ci-runs-full-suite.md)
- decided-by [scoped-tests-serve-m25-generic-rule](../../../decision/scoped-tests-serve-m25-generic-rule.md)

### In

- verifies [guides-scoped-tests-verified](../../../evidence/guides-scoped-tests-verified.md)
- verifies [guides-scoped-tests-qa-1](../../../evidence/guides-scoped-tests-qa-1.md)
- verifies [pr-259-agents-in-runs-stop-running](../../../evidence/pr-259-agents-in-runs-stop-running.md)
- verifies [pr-521-build-role-topics](../../../evidence/pr-521-build-role-topics.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
