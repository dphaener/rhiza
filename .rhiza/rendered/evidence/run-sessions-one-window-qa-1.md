# run-sessions-one-window-qa-1

evidence · `01M66J3DXNGG1BDP1MQACKJ4J9` · actor [claude-code-d1914dd4](../actor/claude-code-d1914dd4.md)

## result

PASS

## statement

Built the branch head (go build -o $TMPDIR/rhiza-x ./cmd/rhiza) and ran rhiza run start --workflow review (checks alpha, beta, gamma, with critic, fixer and verifier) in a scratch repo against the real Rex server, with only claude faked (records its round, says rhiza session done, stays open like a real agent). The run succeeded. The logged rex calls were exactly one window new, then six splits (check beta, check gamma, critic alpha #1, fixer thing #2, verifier thing #2, check alpha #2), each against a live block of that window, then block close of every pane at step end.

## Links

### Out

- verifies [run-sessions-one-window](../capability/rhiza/runtime/rex-sessions/run-sessions-one-window.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
