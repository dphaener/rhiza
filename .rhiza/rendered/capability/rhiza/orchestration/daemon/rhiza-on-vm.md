# rhiza-on-vm

capability · `01M6FESPV5PW1DHSM7HQC0F4CP` · parent [daemon](../daemon.md)

## acceptance

Evidence from the VM: rhiza daemon status --json shows mode service; rhiza run ls shows a ship-change run the VM's daemon started ending done with its PR merged; a question asked in an agent tab on the VM raised a Rex notification on the Mac, or a report was filed with Rex upstream; rhiza daemon status on the Mac shows no daemon running

## statement

The root daemon runs as a service (rhiza-root.service, on the root ~/code where this repo is registered) on the VM and ships a change set end to end with no orchestrator session, worked from the Mac's Rex app connected to the VM host, while the Mac runs no daemon

## Links

### In

- blocks [vm-move-who](../../../../question/vm-move-who.md)
- blocks [vm-move-mac-checks](../../../../question/vm-move-mac-checks.md)
- blocks [leaf-failed-01m6h8hy-1](../../../../question/leaf-failed-01m6h8hy-1.md)
- blocks [leaf-failed-01m6h8hy-2](../../../../question/leaf-failed-01m6h8hy-2.md)
- verifies [rhiza-on-vm-verified](../../../../evidence/rhiza-on-vm-verified.md)
- verifies [rhiza-on-vm-qa-1](../../../../evidence/rhiza-on-vm-qa-1.md)
- verifies [pr-307-move-this-repo-s-rhiza](../../../../evidence/pr-307-move-this-repo-s-rhiza.md)
- verifies [pr-351-this-vm-on-one-root](../../../../evidence/pr-351-this-vm-on-one-root.md)

## History

- `01M6FENG43982SZMA5Q9VBGEVZ` Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook
- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
