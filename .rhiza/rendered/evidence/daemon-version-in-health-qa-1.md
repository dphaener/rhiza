# daemon-version-in-health-qa-1

evidence · `01M6FMSZ9KC7MS0KX2BX4C5E40` · actor [claude-code-c035c61a](../actor/claude-code-c035c61a.md)

## result

PASS

## statement

At f5fda44d I built rhiza with go build from a plain clone (vcs.revision=f5fda44d embedded, no RHIZA_TEST_COMMIT) and ran a real daemon (paused, fake rex) in that clone. rhiza daemon status printed 'Commit: f5fda44df4dfe5683910da14c1a72b9b5aa4328d' and --json carried "commit". With refs/remotes/origin/main at that commit, rhiza health's daemon part was 'ok'. After moving origin/main to a child commit, it was 'warn ... ; built from f5fda44, origin/main is at fdb2362' with next 'make update' (also in --json). A diverged origin/main and no origin/main both gave 'ok'.

## Links

### Out

- verifies [daemon-version-in-health](../capability/rhiza/orchestration/daemon/daemon-version-in-health.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs
