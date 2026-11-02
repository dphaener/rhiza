# close-out-pull-recorded-verified

evidence · `01M6DF2AR5NKR7SS8R2NK81H67` · actor [claude-code-519f11df](../actor/claude-code-519f11df.md)

## result

PASS

## statement

At 82a2fe57 (verifier claude-code-519f11df): `go test ./internal/api -run TestShipCloseOutRecordsStaleMain -count=1 -v` PASS. The test runs ship.close-out with a dirty main checkout (status ' M README.md'). It checks that the inbox lists one main-stale item naming 'has uncommitted changes' and that healthMain is warn, with summary 'the last close-out could not pull' and detail naming the change set and the failure. Then a later close-out pulls, and the test checks that both the inbox item and the health mention are gone. `go test ./cmd/rhiza -run 'TestScripts/(inbox-main-stale|health)$' -count=1` ok, and inbox open now finds the item. In the diff, internal/api/health.go healthMain reads daemon.ReadMainStale, the main part the health op from health-app adds; that op was missing in attempts 1 and 2. `make check` exits 0 with RHIZA_ROOT/CHANGE/RUN unset.

## Links

### Out

- verifies [close-out-pull-recorded](../capability/rhiza/orchestration/daemon/daemon-status/health-command/close-out-pull-recorded.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
