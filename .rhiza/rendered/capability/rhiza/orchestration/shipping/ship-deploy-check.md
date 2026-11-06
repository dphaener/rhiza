# ship-deploy-check

capability · `01M6J3FFWKR8A3D1ZRM6APR30A` · parent [shipping](../shipping.md)

## acceptance

Workflow tests against a local fake health server with no real network: 503, 503, then 200 ends the run done; 503 until the timeout fails the run with the last status in its journal and in rhiza inbox; with no URL set the run ends done without polling.

## statement

After a ship-change run merges, when the repo has a deploy health URL it polls the URL until it answers 200 or the deploy timeout passes, and on timeout fails the run with the last status it saw, which the inbox shows; with no URL the step is skipped.

## Links

### Out

- decided-by [deploy-check-daemon-setting](../../../../decision/deploy-check-daemon-setting.md)
- decided-by [nebari-tooling-retired](../../../../decision/nebari-tooling-retired.md)

### In

- verifies [ship-deploy-check-verified](../../../../evidence/ship-deploy-check-verified.md)
- verifies [ship-deploy-check-qa-1](../../../../evidence/ship-deploy-check-qa-1.md)
- verifies [pr-332-deploy-check-on-ship](../../../../evidence/pr-332-deploy-check-on-ship.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
