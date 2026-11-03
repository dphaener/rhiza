# m24-on-the-vm

intent · `01M6FEM51ZKB5VM5YCY867AR0A` · parent [milestone-24](../milestone-24.md)

## statement

This repo's Rhiza installation runs on a dedicated Ubuntu Server 26 VM on Proxmox, and the owner works it from the Mac's Rex app with the Mac's daemon stopped

## Links

### In

- serves change set `01M6FSQG5W0V2D5PFJT8CHK3DV` (Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them)
- serves change set `01M6GCPSMFW8G6MXXNW6Y5J0Q6` (Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state)
- serves change set `01M6GG7VES74Y7T74DD8HKXFJ7` (Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu)
- serves change set `01M6H6SGKR1N5GBG7XPECKN9Z3` (Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
