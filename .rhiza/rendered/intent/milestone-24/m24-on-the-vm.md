# m24-on-the-vm

intent · `01M6FEM51ZKB5VM5YCY867AR0A` · parent [milestone-24](../milestone-24.md)

## statement

This repo's Rhiza installation runs on a dedicated Ubuntu Server 26 VM on Proxmox, and the owner works it from the Mac's Rex app with the Mac's daemon stopped

## Links

### In

- serves change set `01M6FSQG5W0V2D5PFJT8CHK3DV` (Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
