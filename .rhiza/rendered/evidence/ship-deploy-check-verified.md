# ship-deploy-check-verified

evidence · `01M6RZ2DAPM61CXCC6JEZTD8ZF` · actor [claude-code-b5c2adfd](../actor/claude-code-b5c2adfd.md)

## result

PASS

## statement

At ac624188 I read the diffs of a5a10eeb and ac624188 and ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-deploy-check|run-ship-deploy-timeout|run-ship-change|run-ship-change-rebase|describe-golden)$' -count=1: all PASS. The tests use only a loopback httptest health server and fake rex, claude and gh. run-ship-deploy-check runs the built-in ship-change against a server answering 503, 503, then 200: the run succeeds, deploy-check passes after ship, and the server log is exactly 503/503/200. run-ship-deploy-timeout answers 503 with a 300ms timeout: the run fails at deploy-check, the journal's step and finished entries hold 'still answered HTTP 503', and rhiza inbox, both human and --json, shows that status. run-ship-change sets no URL: it prints the skip line, never logs a poll, and deploy-check passes. go test ./internal/api -run ShipDeploy and ./internal/workflow -run ShipChange also pass, and gofmt -l and go vet are clean on internal/api, cmd/rhiza and internal/workflow.

## Links

### Out

- verifies [ship-deploy-check](../capability/rhiza/orchestration/shipping/ship-deploy-check.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
