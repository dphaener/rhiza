# rhiza-on-vm-wait-for-mac-checks

decision · `01M6H9STBS5P00PE10GPD3PZ5D`

## reasoning

All three attempts failed on the same two acceptance points, and both can only be read on the owner's Mac: commits 48ceb789 and c4e83e8e show the VM cannot reach it (ssh refused, Rex host unreachable over Tailscale). The VM half is done and recorded on vm-move-mac-checks (daemon in service mode under systemd; ship-change run 01M6GPD3P4SD started by the VM daemon succeeded with PR #212 merged). No builder can produce the Mac evidence, so another round would fail the same way. Dropping or weakening the two points would hollow the capability's claim that the Mac runs no daemon and that notifications reach the owner, which is the point of the move, so the acceptance stays. vm-move-mac-checks already blocks the leaf and asks the owner exactly what is missing; the right path is to let it gate the leaf rather than spend attempts.

## statement

Keep rhiza-on-vm and its acceptance as written, and stop rebuilding it: the leaf waits on question vm-move-mac-checks. Once the owner answers that question with the Mac's Rex notification (or an upstream Rex report) and the Mac's rhiza daemon status, the next attempt records those answers as Evidence beside the VM-side Evidence already gathered, with no code change.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6h8hy-1](../question/leaf-failed-01m6h8hy-1.md)

## History

- `01M6FENG43982SZMA5Q9VBGEVZ` Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook
