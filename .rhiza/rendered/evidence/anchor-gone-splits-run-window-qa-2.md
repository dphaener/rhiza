# anchor-gone-splits-run-window-qa-2

evidence · `01M66KA1R7Z7D02PH0SWB3KD4J` · actor [claude-code-cec4c5b5](../actor/claude-code-cec4c5b5.md)

## result

PASS

## statement

At 68bf3aa, built rhiza (go build) and ran rhiza run start against the real Rex server in scratch repos, with only claude faked. Review run where check alpha's and the critic's agents exit (panes gone) while beta stays open: the fixer's rex split -b <critic's closed block> failed not found, rex block ls -w <run window> listed beta's live block, and the fixer ran rex split -b <beta's block> in the same window; the whole run made one rex window new and succeeded. Two-step workflow (the cycle-1 failure): step one's panes closed at step end, step two's split of the gone anchor failed not found, rex block ls reported the window not found, and the session opened a new window; the run succeeded. Review run with every agent exiting: each session whose window was gone opened a new window and the run succeeded; no Panes windows were left over.

## Links

### Out

- verifies [anchor-gone-splits-run-window](../capability/rhiza/runtime/rex-sessions/anchor-gone-splits-run-window.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
