# maxruns-auto-setting-verified

evidence · `01M7QXW25R2YHA603M9ZQQPHJV` · actor [claude-code-646cf1a1](../actor/claude-code-646cf1a1.md)

## result

PASS

## statement

In a detached checkout at 0b9e4a2b9: go test -trimpath -p 2 ./internal/api -run 'MaxRuns' exit 0 (TestDaemonMaxRunsAuto PASS: daemon set auto, daemon settings reads auto, daemonConf maxRunsAuto true with maxRuns at the default; 3 pins cap 3; --max-runs 4 overrides auto; 0, -1, autox refused naming 'a whole number of at least 1 or auto'; auto refused for rhiza.project.maxRuns); go test -trimpath -p 2 ./internal/tui/daemonapp/... exit 0 (TestAutoLimitSetsAndQuits writes daemon.set rhiza.daemon.maxRuns auto; TestLimitValidators, TestInvalidLimitCallsNothing refuse 0/-1/autox); -count=1 rerun of internal/api -run 'MaxRuns|Setting|DaemonConf' exit 0; rhiza ship generate and rhiza render exit 0 with a clean tree. Diff read against c8e833d39: default 2, --max-runs int and rhiza.project.maxRuns unchanged, poll untouched, no stubs or TODOs; only the leaf's files plus the generated docs/reference/config.md and .rhiza changed.

## Links

### Out

- verifies [maxruns-auto-setting](../capability/rhiza/orchestration/daemon/daemon-settings/maxruns-auto-setting.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
