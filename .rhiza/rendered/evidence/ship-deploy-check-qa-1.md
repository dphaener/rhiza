# ship-deploy-check-qa-1

evidence · `01M6S15KAGKCZ33SR1YHKWKB3H` · actor [claude-code-4f092051](../actor/claude-code-4f092051.md)

## result

PASS

## statement

QA at 5872d5ea with a scratch build against real local python HTTP health servers on loopback: with no URL set, rhiza ship deploy prints the skip line and exits 0 (--json skipped:true, polls:0). Against a server answering 503, 503, then 200 (RHIZA_DEPLOY_INTERVAL=200ms, timeout 10s) it printed deploy check 1/2: HTTP 503, then Deploy healthy after 3 polls, exit 0; the server logged exactly 503/503/200. Against a server always answering 503 with timeout 1s it polled 4 times and printed Deploy not healthy ... still answered HTTP 503 after 4 polls over 1s, exit 1 (--json last: HTTP 503). With nothing listening it fails exit 1 with the connection error as the last status. rhiza api call ship.deploy reaches the same handler. The whole ship-change run (merge, then deploy-check failing the run, journal and rhiza inbox showing HTTP 503; no-URL run skipping) was exercised through the scoped testscripts run-ship-deploy-check, run-ship-deploy-timeout and run-ship-change, all passing, since a real merge needs gh; go test ./internal/api -run ShipDeploy passes.

## Links

### Out

- verifies [ship-deploy-check](../capability/rhiza/orchestration/shipping/ship-deploy-check.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
