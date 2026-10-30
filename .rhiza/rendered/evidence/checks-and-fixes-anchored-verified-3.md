# checks-and-fixes-anchored-verified-3

evidence · `01M66M1DKQM36SCRNYC2BEFX45` · actor [claude-code-fb34d4b4](../actor/claude-code-fb34d4b4.md)

## result

PASS

## statement

At d421578 ran go test ./internal/api -run 'TestStartCheck|TestStartFix' -count=1 -v: TestStartCheck, TestStartFix and TestStartCheckLatestAnchor all PASS. Read the tests: each runs StartCheck/StartFix with an Anchor (S-first, block:a1) and wantSplitOf asserts a rex split -b block:a1 and fails on any rex window new; TestStartCheckLatestAnchor shows a session with no Anchor splits the run's latest pane (block:a2). paneAnchor (internal/api/run_leaves.go:120) is used by StartCheck and StartFix (run_checks.go:59,92). make check passed fmt-check, vet and test; self-check failed only on the render of this verifier's own actor record, not on the commit.

## Links

### Out

- verifies [checks-and-fixes-anchored](../capability/rhiza/runtime/rex-sessions/checks-and-fixes-anchored.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
