# build-leaf-fails-verified

evidence · `01M5ZMT73521QGV6MJ84APMRQ1` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestLeafFails, over leaves a, b and c with c depending on a and a failing: one question raised for a (AddLeafQuestion, which TestAddLeafQuestion shows writes a question linked blocks to the leaf), c skipped, b built and verified, the run failed. TestLeafFailsResume: after an answering Decision, Retry starts a's builder afresh (attempt 4, round 2, ladder's first model) with the Decision and its answer in its spec, then c; testscript build-leaf-fails finds the Decision text in the real builder's prompt and order alpha 4 then beta 1. TestLeafFailsAutonomous: the arbiter is started on a's question and a second failure ends the run failed. All passed.

## tests

go test -count=1 ./internal/run -run TestLeafFails; ./internal/api -run 'TestAddLeafQuestion|TestLeafPrompt'; ./cmd/rhiza -run 'TestScripts/build-leaf-fails$' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-leaf-fails](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-fails.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
