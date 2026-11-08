# ship-head-build-error-detail

capability · `01M6XHVSXCAK8XB69G5FFYN9MT` · parent [shipping](../shipping.md)

## acceptance

A Go test in internal/api with a fake Exec whose head build exits 3 printing {"error":{"code":"conflict","message":"cannot land change set X: ..."}} shows the ship land error naming the conflict message and exit status 3; a second case with no JSON body shows the child's stderr in the error; run with go test -trimpath ./internal/api/... -run 'Ship'

## statement

When a ship step runs rhiza built from the PR head (land, render) and that run fails, the step's error carries the child's error message and code from its --json error body (or its stderr when there is none) besides the exit status, so the run log says why, e.g. the land conflict, instead of only 'exit status 3'

## Links

### Out

- derived-from [plan-fix-watchdog-unmerged-pr-at](../../../../source/plan-fix-watchdog-unmerged-pr-at.md)

### In

- verifies [ship-head-build-error-detail-verified](../../../../evidence/ship-head-build-error-detail-verified.md)
- verifies [ship-head-build-error-detail-qa-1](../../../../evidence/ship-head-build-error-detail-qa-1.md)
- verifies [pr-389-fix-watchdog-unmerged-pr-at](../../../../evidence/pr-389-fix-watchdog-unmerged-pr-at.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
