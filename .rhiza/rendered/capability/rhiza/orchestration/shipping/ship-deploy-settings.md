# ship-deploy-settings

capability · `01M6J3F8SSE7GHCRS47KVY29CZ` · parent [shipping](../shipping.md)

## acceptance

A testscript sets both, sees them in rhiza daemon settings --json with source git config, and gets exit 2 with a next for a URL that is not http or https and for a timeout that is not a duration.

## statement

rhiza daemon set takes two per-repo ship settings, a deploy health URL (http or https) and a deploy timeout (a duration, default 15m), and rhiza daemon settings lists both with their value and source; an unset health URL means no deploy check.

## Links

### Out

- decided-by [deploy-check-daemon-setting](../../../../decision/deploy-check-daemon-setting.md)

### In

- verifies [ship-deploy-settings-verified](../../../../evidence/ship-deploy-settings-verified.md)
- verifies [ship-deploy-settings-qa-1](../../../../evidence/ship-deploy-settings-qa-1.md)
- verifies [pr-332-deploy-check-on-ship](../../../../evidence/pr-332-deploy-check-on-ship.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
