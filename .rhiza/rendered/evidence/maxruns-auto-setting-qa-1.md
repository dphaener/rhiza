# maxruns-auto-setting-qa-1

evidence · `01M7QZEJ3E6BH8CP5SJD8Z2YAT` · actor [claude-code-34e635d3](../actor/claude-code-34e635d3.md)

## result

PASS

## statement

QA by hand with a fresh build at 009d74d2c (go build -trimpath to $TMPDIR/rhiza-qa) in a scratch root: rhiza daemon set rhiza.daemon.maxRuns auto exits 0 and daemon settings shows 'auto (rhiza.toml, machine)'; 0, -1, autox, AUTO, 1.5 and empty are refused exit 2 with 'must be a whole number of at least 1 or auto'; rhiza.project.maxRuns auto is refused ('must be a whole number of at least 1'); a hand-edited rhiza.toml with autox makes daemon settings fail naming the same rule; the daemon app's run-limit field (not drivable by pty) passes go test ./internal/tui/daemonapp/... (TestAutoLimitSetsAndQuits, TestLimitValidators) and go test ./internal/api -run MaxRuns passes.

## Links

### Out

- verifies [maxruns-auto-setting](../capability/rhiza/orchestration/daemon/daemon-settings/maxruns-auto-setting.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
