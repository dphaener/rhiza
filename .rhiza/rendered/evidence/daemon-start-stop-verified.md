# daemon-start-stop-verified

evidence · `01M6FMXF0N7BQ7QKS4WPZ4MZTY` · actor [claude-code-8f631d1f](../actor/claude-code-8f631d1f.md)

## result

PASS

## statement

At c5d1b664 I read the diff (internal/api/daemon_startstop.go, daemon.go Mode field, autonomyapp using daemon.start/daemon.stop) and ran go test ./cmd/rhiza -run 'TestScripts/(daemon-start-stop|daemon-install|describe-golden)$' -count=3 (pass) and -v (daemon-start-stop PASS, no leftover 'daemon --autonomous' process afterwards). The script, with fake systemctl, loginctl and rex on PATH: with no unit, status --json says mode tab, start logs 'rex window new -s rhiza-repo rhiza daemon -c ... --keep-open --json -- ... daemon --autonomous', the daemon comes up running, a second start does nothing, stop ends it (running false) and a further stop says no daemon runs; after daemon install, start and stop log exactly 'systemctl --user start rhiza-repo.service' and 'systemctl --user stop rhiza-repo.service' and status says mode service; an owner-written unit leaves mode tab. go test ./internal/tui/autonomyapp -v: TestStartDaemonWhenNotRunning (d then y calls daemon.start on the fake runner when stopped) and TestStopRunningDaemon (d then y calls daemon.stop when running) pass. make check (RHIZA_ROOT/CHANGE/RUN unset) exits 0.

## Links

### Out

- verifies [daemon-start-stop](../capability/rhiza/orchestration/daemon/daemon-start-stop.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
