# rhiza-on-vm-qa-1

evidence · `01M6P93YCWQFE5SS86EWPRQPPD` · actor [claude-code-4e88d302](../actor/claude-code-4e88d302.md)

## result

PASS

## statement

QA on the VM (hostname rhiza) on 2026-11-05 with a binary built from the branch head (rhiza v0.0.0-...-e3a6f705), run from ~/code/rhiza with RHIZA_ROOT/CHANGE/RUN unset. (1) `rhiza daemon status --json`: mode service, running true, process rhiza:1348, last_started 2026-11-05T12:06:16Z for change 01M6NR6EK065; `systemctl --user is-active rhiza-rhiza.service` reports active. (2) `rhiza run ls`: ship-change run 01M6NSMNAFQ3N2TF0V7EDVTZKP (The daemon keeps the Go build) started 12:06:22Z, matching the daemon's last_started, ended succeeded; `gh pr list --state merged` shows #303 'The daemon keeps the Go build' merged 2026-11-05T13:04:27Z; `gh pr view 212` (Machine trust) MERGED 2026-11-03T17:12:23Z. The only orchestrator process on the VM is the passive ~/orchestrator/watch.sh status watcher. (3) `rex client ls` on the VM shows the Mac's Rex app (Rex Beta 0.1.0 (1758), darins-macbook-pro) connected over Tailscale. The Mac notification and the Mac's daemon status cannot be read from the VM; they rest on the owner's answers vm-move-notify-chime-and-marker (chime plus blocked marker on the Mac, live test 2026-11-05) and vm-move-mac-no-daemon ('There is no daemon running on the mac'). Read docs/rex-upstream/notification-banner.md: the banner request is written.

## Links

### Out

- verifies [rhiza-on-vm](../capability/rhiza/orchestration/daemon/rhiza-on-vm.md)

## History

- `01M6FENG43982SZMA5Q9VBGEVZ` Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook
