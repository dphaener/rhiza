# maxruns-auto-setting

capability · `01M7QXK2Z9E3VYWJRD30PF1514` · parent [daemon-settings](../daemon-settings.md)

## acceptance

go test -trimpath ./internal/api -run 'MaxRuns' and go test -trimpath ./internal/tui/daemonapp/...: a test sets rhiza.daemon.maxRuns to auto and reads auto back with daemonConf marked auto, a test sets 3 and gets a pinned cap of 3, and the daemon app accepts auto in its run-limit field; (failure case) 0, -1 and autox are refused for rhiza.daemon.maxRuns with an error naming 'a whole number of at least 1 or auto', and auto is refused for rhiza.project.maxRuns

## files

internal/api/daemon_settings.go, internal/api/daemon_settings_test.go, internal/tui/daemonapp/daemonapp.go, internal/tui/daemonapp/daemonapp_test.go

## must_not

Change the default of 2, the --max-runs flag's whole-number type, or rhiza.project.maxRuns; change how the poll starts runs (leaf maxruns-auto-poll)

## statement

rhiza.daemon.maxRuns accepts auto beside a whole number of at least 1 (rhiza daemon set, rhiza daemon settings, rhiza.toml and the Rex daemon app's run-limit field), its type and description say so, and the daemon's resolved settings (daemonConf) carry whether it is auto

## Links

### Out

- decided-by [maxruns-auto-scope](../../../../../decision/maxruns-auto-scope.md)

### In

- depends-on [maxruns-auto-poll](../daemon-picks/maxruns-auto-poll.md)
- verifies [maxruns-auto-setting-verified](../../../../../evidence/maxruns-auto-setting-verified.md)
- verifies [maxruns-auto-setting-qa-1](../../../../../evidence/maxruns-auto-setting-qa-1.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
