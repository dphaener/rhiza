# ship-head-build-error-detail-qa-1

evidence · `01M6XT12952XCFR0YRF91R0SVA` · actor [claude-code-2d293fb3](../actor/claude-code-2d293fb3.md)

## result

PASS

## statement

QA at 27e17a09: produced a real land conflict with the built binary (old binary wrote update + retire on the landed base, new binary ran land --json: exit 3, error body on stdout, nothing on stderr). A throwaway in-package test ran that real binary through os/exec as cmd/rhiza execIn does and passed its stdout and error to headFailure, which returned 'conflict: cannot land change set ...: model: conflict: base ... is stale ...: rhiza land: exit status 3', naming the code, the conflict message and the exit status. Read execIn: it returns stdout with the error, so shipRunHead sees the body. go test -trimpath ./internal/api -run Ship ok.

## Links

### Out

- verifies [ship-head-build-error-detail](../capability/rhiza/orchestration/shipping/ship-head-build-error-detail.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
