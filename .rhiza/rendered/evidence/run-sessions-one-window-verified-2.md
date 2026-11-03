# run-sessions-one-window-verified-2

evidence · `01M66JK0JZ9TMJJWE9YJK385SZ` · actor [claude-code-3597568f](../actor/claude-code-3597568f.md)

## result

PASS

## statement

At 5b41015, go test ./cmd/rhiza -run 'TestScripts/run-panes-one-window$' -count=1 -v passed. With a fake rex and fake claude, a review run started check sessions alpha, beta and gamma, a critic for alpha, then in cycle 2 the fixer and verifier of leaf thing and alpha again (claude.log matched want.log). rex.log held exactly one 'window new' followed by six 'split block:bN' calls, each against a block the fake rex returned in window:w1. Read the diff: StartLeaf, StartCheck and StartFix now take req.anchor from paneAnchor, which falls back to the tab of the run's first session with a pane. go test ./internal/api -count=1 passed. The open review findings (qa f1/f2: a gone run window fails the next session) are about anchor-gone handling, not this leaf's acceptance.

## Links

### Out

- verifies run-sessions-one-window (retired)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
