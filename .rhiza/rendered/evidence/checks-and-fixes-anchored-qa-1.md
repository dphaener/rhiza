# checks-and-fixes-anchored-qa-1

evidence · `01M66J3ESBFRHZE8A6S0ZRZWCP` · actor [claude-code-d1914dd4](../actor/claude-code-d1914dd4.md)

## result

PASS

## statement

In the same real-Rex run (built binary, fake claude only), the critic, the change-set leaf fixer and verifier, and the cycle-2 check each ran rex split -b <block of a live pane of the run's window> and none ran rex window new; each fake agent logged a REX_BLOCK inside that window.

## Links

### Out

- verifies checks-and-fixes-anchored (retired)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
