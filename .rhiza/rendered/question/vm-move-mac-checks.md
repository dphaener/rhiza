# vm-move-mac-checks

question · `01M6H8Q8CWMEZXV640MZ8VEP4E`

## statement

The rhiza-on-vm leaf needs two checks only the owner's Mac can give, and the VM cannot reach the Mac (ssh refused, its Rex host unreachable over Tailscale): did a question asked in an agent tab on the VM raise a Rex notification on the Mac (or was a report filed with Rex upstream), and does `rhiza daemon status` on the Mac show no daemon running? From the VM on 2026-11-03: `rhiza daemon status --json` shows mode service, running, under systemd user unit rhiza-rhiza.service (active); `rhiza run ls` shows ship-change run 01M6GPD3P4SDTAB9X4778TJJE4 (Machine trust), started by the VM daemon after the 05:45Z handoff, succeeded, its PR #212 merged 2026-11-03T17:12:23Z; ship-change run 01M6G29SNV0B6V7SZHHYNC7K3Z (Fix) also succeeded.

## Links

### Out

- blocks [rhiza-on-vm](../capability/rhiza/orchestration/daemon/rhiza-on-vm.md)

### In

- answers [vm-move-mac-no-daemon](../decision/vm-move-mac-no-daemon.md)
- answers [vm-move-notify-chime-and-marker](../decision/vm-move-notify-chime-and-marker.md)
- discusses [note-01m6h8ya0f8vwn](../note/note-01m6h8ya0f8vwn.md)
- discusses [note-01m6h9zyb93sb7](../note/note-01m6h9zyb93sb7.md)
- discusses [note-01m6ha3hvkw24y](../note/note-01m6ha3hvkw24y.md)

## History

- `01M6FENG43982SZMA5Q9VBGEVZ` Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook
