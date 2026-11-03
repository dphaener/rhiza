# daemon-start-stop-qa-1

evidence · `01M6FN79S4J7TT91TXFXM266A7` · actor [claude-code-e6cae1c1](../actor/claude-code-e6cae1c1.md)

## result

PASS

## statement

QA at 78a3b478 with a binary built from the branch, in scratch repo qarepo on the real Rex server. Tab mode: daemon status said 'Mode: tab'. daemon start --json returned {changed:true, mode:tab}, opened the Rex tab 'rhiza daemon' in session rhiza-qarepo, and status then said 'Daemon running'. A second start said 'The daemon already runs' and opened no tab. daemon stop said 'Asked the daemon to stop'; the process exited and status said not running. A further stop said 'No daemon runs'. Service mode, after daemon install with fake systemctl/loginctl and a scratch HOME: status said 'Mode: service', and start and stop logged exactly 'systemctl --user start rhiza-qarepo.service' and 'systemctl --user stop rhiza-qarepo.service'. After uninstall the mode was tab again. On a fresh repo with no permission mode set, start reported success, but the tab's daemon --autonomous exited at once ('needs --permission-mode'). See the minor finding.

## Links

### Out

- verifies [daemon-start-stop](../capability/rhiza/orchestration/daemon/daemon-start-stop.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
