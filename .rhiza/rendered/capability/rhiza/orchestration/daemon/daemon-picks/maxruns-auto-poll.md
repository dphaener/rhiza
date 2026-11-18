# maxruns-auto-poll

capability · `01M7QXK319H58FEFW7GQ2B3ZWY` · parent [daemon-picks](../daemon-picks.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/daemon-maxruns-auto: a new script with three accepted change sets, rhiza.daemon.maxRuns=auto and idle readings starts one run on the first --once poll and a second on the next (readings with room start a second run), and the saved state has max_runs_auto true; with maxRuns set to 1 a further poll starts none; (failure case) with RHIZA_TEST_LOAD1 and RHIZA_TEST_LOAD5 over rhiza.daemon.maxLoad, and again with the tmp-low hold on, an auto poll starts none; go test -trimpath ./internal/api/... and ./internal/daemon/... pass

## files

internal/api/daemon.go, internal/api/daemon_settings.go, internal/daemon/state.go, cmd/rhiza/testdata/script/daemon-maxruns-auto.txtar

## must_not

Read the host's load, memory or free space anew for the cap (use the poll's own pre-start readings); start more than one run per poll in auto; change a numeric maxRuns' behaviour; change rhiza run spawn's holds

## statement

With rhiza.daemon.maxRuns=auto the root daemon's poll sets the machine-wide cap it fills round-robin to the runs running plus one while every pre-start check has room, and to the runs running while any holds; the saved daemon state records max_runs_auto true beside max_runs, the cap that poll used; a number set for maxRuns pins the cap as before

## Links

### Out

- depends-on [maxruns-auto-setting](../daemon-settings/maxruns-auto-setting.md)
- decided-by [maxruns-auto-one-per-poll](../../../../../decision/maxruns-auto-one-per-poll.md)

### In

- depends-on [maxruns-auto-shown](../daemon-status/maxruns-auto-shown.md)
- verifies [maxruns-auto-poll-verified](../../../../../evidence/maxruns-auto-poll-verified.md)
- verifies [maxruns-auto-poll-qa-1](../../../../../evidence/maxruns-auto-poll-qa-1.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
