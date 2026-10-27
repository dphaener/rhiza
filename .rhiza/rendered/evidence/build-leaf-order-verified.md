# build-leaf-order-verified

evidence · `01M5ZMT56VESAHFQGD1C257BNX` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestLeafOrderOverlapAndGating runs a fake runner over leaves a, b and c, c depending on a, with a's verifier live for five polls: it finds b's builder session started after a's verifier started and before a's verifier was judged, and c's builder not started before 'leaf a passed' (which follows a's Evidence check and commit), with only one builder live at a time. TestLeafOrderNextActions and TestLeafOrderResume cover the same gating on the next-action level and on a resume. All passed.

## tests

go test -count=1 ./internal/run -run 'TestLeafOrderOverlapAndGating|TestLeafOrderNextActions|TestLeafOrderResume' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-leaf-order](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-order.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
