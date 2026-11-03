# milestone-24

intent · `01M6FEKX9J8ARG7JGDYSZEA2YW`

## milestone

24

## statement

Rhiza runs headless on a dedicated Linux VM: its daemon runs as a service beside Rex's Linux server, drives runs with no tabs of its own, keeps itself current with main, and the owner works it from the Mac's Rex app, where only agent work shows

## Links

### Out

- decided-by [m24-server-daemon-no-tab](../decision/m24-server-daemon-no-tab.md)
- decided-by [m24-only-agent-tabs](../decision/m24-only-agent-tabs.md)

### In

- serves change set `01M6FEMDQ1D2HWXT3GR4F3Z5H6` (The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session)
- serves change set `01M6FEN1BRQRHP9RTJKN2YPTH6` (The daemon pulls main and runs make update when main moves, and health shows the version it runs)
- serves change set `01M6FENEJE4B11EBYXDGP9T052` (Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs)
- serves change set `01M6FENG43982SZMA5Q9VBGEVZ` (Move this repo's Rhiza installation to a dedicated Ubuntu VM on Proxmox, by an HTML runbook)
- serves change set `01M6GDJW3R8PXW7QBAZVH0QS4W` (Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix')
- serves change set `01M6GJ9GBCQD9GFCDSNR4BB13Q` (Serialize the three bug change sets titled Fix, so their ship runs never share cs/fix and rhiza.fix: the test-scripts fix waits for the resume fix)
- serves change set `01M6H3TC786A6WVKRXC0WQ1N3K` (Agent panes close when their session finishes: every driver closes a finished session's pane at once (review checks, critic, fixers, plan steps), not at close-out)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
