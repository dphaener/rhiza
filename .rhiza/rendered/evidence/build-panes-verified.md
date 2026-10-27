# build-panes-verified

evidence · `01M5ZMT7D9TQBJGEY8ERWAMNND` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Fake rex: TestOpenSessionPane finds split right for a landscape size, below for a portrait one, and the override honored without asking the size; TestSplitFor covers the same. Testscript build-panes with a fake rex finds, landscape, every split right, b's verifier splitting a's live verifier pane (the verifier half), a passed verifier's pane closed and the failed sessions' panes closed only at the step's end; portrait splits below; split = "below" overrides a landscape pane. TestLeavesStepPanes finds the same anchors and close order in the journal. All passed.

## tests

go test -count=1 ./internal/session -run TestSplitFor; ./internal/api -run 'TestOpenSessionPane|TestCloseSessionPane'; ./internal/run -run 'TestLeafPaneAnchor|TestLeavesStepPane|TestLeavesStepLivePaneStays'; ./cmd/rhiza -run 'TestScripts/build-panes$' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-panes](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-panes.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
