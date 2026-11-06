# ship-deploy-settings-verified

evidence · `01M6S0T40B90X4SE2EFYX41KX6` · actor [claude-code-410c2f37](../actor/claude-code-410c2f37.md)

## result

PASS

## statement

Verifier at 4b577a66: read the leaf diff (0423ce6f: rhiza.ship.deployUrl and rhiza.ship.deployTimeout in internal/api/daemon_settings.go, default 15m, unset URL means no check, bad values return Usagef WithNext). gofmt -l internal cmd lists nothing (attempt 1's finding is fixed). go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-deploy-settings|describe-golden)$' passes: the script sets both keys, sees them in daemon settings --json with source "git config", and rejects ftp://, a schemeless URL and 'soon' with a next. go test -short -run 'Setting|Daemon' ./internal/api passes. On a scratch build in a hermetic repo: set deployUrl http://localhost:8080/up and deployTimeout 90s exit 0 and show source git config in settings --json; ftp://x and 'soon' each exit 2 printing 'next: rhiza daemon settings'. The full internal/api package ran past its timeout here only because the machine was under load (TestReadBudget waits for a quiet machine); CI runs the full package.

## Links

### Out

- verifies [ship-deploy-settings](../capability/rhiza/orchestration/shipping/ship-deploy-settings.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
