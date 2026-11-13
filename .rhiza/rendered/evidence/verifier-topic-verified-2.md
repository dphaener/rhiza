# verifier-topic-verified-2

evidence · `01M79BJSF7GM4BM0FGCT3F0K5W` · actor [claude-code-189a522c](../actor/claude-code-189a522c.md)

## result

PASS

## statement

Verifier attempt 2 at 8dfc21dab: go test -trimpath ./internal/guide/... exit 0; TestVerifierText passes and now asserts the exact six ## headings (Task, Expected outcome, Read first, Must do, Must not do, Context) with command-budget nested in Must do and never-list nested in Must not do, plus session-done/test-scope blocks, detached checkout, Never edit code, fakes only, Checks:/acceptance re-run with missing block or different exit code a finding naming the command, scope and no-stub points, the Evidence apply example, the rhiza run findings heredoc, findings naming file/command/acceptance point, and the negative case banning builder commit/trailer instructions. A scratch build's `rhiza guide verifier` prints only those six ## sections (### Command budget, ### Never nested). Re-ran every Checks: command: the -run TestBuilderText|TestVerifierText|TestReviewText|TestShared|TestBuildReviewCommandBudget run, go test ./cmd/rhiza -run TestClaudeMDScopedTests, go test ./internal/workflow, rhiza ship generate, rhiza render (tree stayed clean), make self-check: all exit 0. Diff touches only verifier.md, verifier_test.go and .rhiza records; the seven stale guide nodes (verifier-reruns-reported-checks, unstable-ci-guidance, delay-repro-guidance, guides-scoped-tests, command-budget-guide, build-guide-fakes-only, done-own-command-guide) no longer name build.md or TestBuildText. Review findings f1 (quality, qa) and f3 (model) are addressed.

## Links

### Out

- verifies [verifier-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-topic.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
