# guide-shared-blocks-verified

evidence · `01M799CCKHRQA53DT1KHSWV363` · actor [claude-code-43876667](../actor/claude-code-43876667.md)

## result

PASS

## statement

Verifier at eaa71074c (attempt 2): go test -trimpath -count=1 ./internal/guide/... passes (TestShared*, TestBuild*, TestBuildReviewCommandBudget all PASS). Read the diff: guide.go embeds text/shared/*.md and Expand replaces each `<!-- shared: <block> -->` line at load, panicking with the file name on an unknown block. shared_test.go covers (1) every shared block expands in place and no listed text keeps a marker, (2) Expand of no-such-block errors naming it, (3) every List() text within its cap via checkCap, (4) never-list forbids sleep/poll, stubs/TODOs/placeholders, new dependency, rewriting history and chaining session done, and every budget, scoped-tests, unstable-ci, outside-work and done phrase of build_test.go and budget_test.go (now including rhiza change new, discovered-from <this>, rhiza change use <this>, the attempt-1 finding) is in exactly one shared file. build.md uses six markers; a scratch build's `rhiza guide build` prints no marker.

## Links

### Out

- verifies [guide-shared-blocks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/guide-shared-blocks.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
