# daemon-self-update-qa-1

evidence · `01M6FMT1BMA5TGHGNW8Y857DCH` · actor [claude-code-c035c61a](../actor/claude-code-c035c61a.md)

## result

PASS

## statement

At f5fda44d I ran the built binary as a real daemon (paused, selfUpdate set with rhiza daemon set; fake rex) in a clone of a scratch bare origin whose Makefile's update target logs and touches the daemon's binary. rhiza daemon settings first listed 'rhiza.daemon.selfUpdate = false (default)', and 'maybe' was refused. Pushing v2 to origin: on its first poll the daemon fetched, fast-forwarded main to 7f19db9, ran make update once, logged 'self-update: updated to 7f19db900', then 'binary changed, re-executing'; two unmoved polls ran nothing. A dirty checkout with origin moved logged 'skipped ... uncommitted changes' once and was not retried after cleaning until origin moved again, when it updated. A failing make update (exit 2) logged once and was not retried. A checkout on topic and a diverged main were each logged and skipped. The pending-start gate was not exercised by hand.

## Links

### Out

- verifies [daemon-self-update](../capability/rhiza/orchestration/daemon/daemon-self-update.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs
