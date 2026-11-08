# ship-head-build-error-detail-verified

evidence · `01M6XPQAXVGX48H4KNAFEGGDB2` · actor [claude-code-7e796d6a](../actor/claude-code-7e796d6a.md)

## result

PASS

## statement

Verifier at 737fcfc5: ran go test -trimpath ./internal/api/... -run 'Ship' (ok) and TestShipLandHeadBuildErrorDetail -v: case 'json body' fakes the head build bin exiting with err 'exit status 3' and stdout {"error":{"code":"conflict","message":"cannot land change set X: node N changed"}}; ship.land's error names conflict, the message and exit status 3 (PASS). Case 'stderr only' with no JSON body and an err carrying the child's stderr (as cmd/rhiza execIn builds it) shows that stderr in the error (PASS). Diff read: shipRunHead now formats headFailure(out, err), parsing the --json error body from stdout for land and render; gofmt and go vet clean.

## Links

### Out

- verifies [ship-head-build-error-detail](../capability/rhiza/orchestration/shipping/ship-head-build-error-detail.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
