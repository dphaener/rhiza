# ship-deploy-settings-qa-1

evidence · `01M6S15K2JZTAM1325Q036TX6G` · actor [claude-code-4f092051](../actor/claude-code-4f092051.md)

## result

PASS

## statement

QA at 5872d5ea with a scratch build (go build ./cmd/rhiza) in a hermetic git repo after rhiza init: rhiza daemon settings shows rhiza.ship.deployUrl = "" (default) and rhiza.ship.deployTimeout = 15m (default). daemon set of deployUrl http://127.0.0.1:18931/health and deployTimeout 10s exit 0, and settings --json lists both with source "git config" and their defaults. ftp://x, example.com (no scheme), soon, -5m and 0s each exit 2 with a message (http or https URL / positive duration) and next: rhiza daemon settings. Scoped testscripts ship-deploy-settings and describe-golden pass.

## Links

### Out

- verifies [ship-deploy-settings](../capability/rhiza/orchestration/shipping/ship-deploy-settings.md)

## History

- `01M6J375CH04RMRS4NN372NCN2` Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza
