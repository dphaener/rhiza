# rhiza-on-vm-verified

evidence · `01M6P904CMYNP9RNRD6MABG6VV` · actor [claude-code-f7a1f570](../actor/claude-code-f7a1f570.md)

## result

PASS

## statement

Verified at e71d9819 on the VM on 2026-11-05. (1) `rhiza daemon status --json` reports mode service, running, process rhiza:1348; `systemctl --user is-active rhiza-rhiza.service` reports active. (2) `rhiza run ls` lists ship-change runs the VM daemon started ending succeeded, e.g. 01M6NSMNAFQ3N2TF0V7EDVTZKP (The daemon keeps the Go build, started 2026-11-05T12:06Z, matching the daemon's last_started), whose PR #303 is merged per `gh pr list --state merged`; the earlier VM-started run 01M6GPD3P4SD (Machine trust) has PR #212 MERGED 2026-11-03T17:12:23Z per `gh pr view 212`. (3) Notification: question vm-move-mac-checks is answered by vm-move-notify-chime-and-marker, recording the owner's live test on 2026-11-05 in which blocked reports from VM agent tabs chimed on the Mac (tests C to G) and one showed a banner (test B); the owner accepted chime plus blocked marker as the notification, and the banner request for Rex upstream is written at docs/rex-upstream/notification-banner.md. (4) Mac daemon: answer vm-move-mac-no-daemon quotes the owner: 'There is no daemon running on the mac'. Read the diffs of dc9ccaf6 and e71d9819: model records and the upstream request doc only, no code.

## Links

### Out

- verifies [rhiza-on-vm](../capability/rhiza/orchestration/daemon/rhiza-on-vm.md)

## History

- `01M6FENG43982SZMA5Q9VBGEVZ` Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook
