# run-sessions-one-window-verified

evidence · `01M66HKBGPXD2HJV0K6WN2NGER` · actor [claude-code-8f1e24fd](../actor/claude-code-8f1e24fd.md)

## result

PASS

## statement

At e28b6c7 I ran go test ./cmd/rhiza -run 'TestScripts/run-panes-one-window$' -count=1 -v: PASS. With the fake rex and claude, the run started checks alpha, beta and gamma, alpha's critic, thing's fixer and verifier, and alpha again in cycle 2 (claude.log matched want.log). rex.log held exactly one 'window new' followed by six 'split block:bN' lines against blocks of window w1 (b1..b5), so the review, critic and fixer sessions all split the run's one window. I read paneAnchor in internal/api/run_leaves.go, used by StartLeaf, StartCheck and StartFix: it falls back to the tab of the run's first session with a block, so only the first session opens a window. make fmt-check vet test pass; make self-check fails only on rendered actor files written during this run, not on the leaf's code.

## Links

### Out

- verifies [run-sessions-one-window](../capability/rhiza/runtime/rex-sessions/run-sessions-one-window.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
