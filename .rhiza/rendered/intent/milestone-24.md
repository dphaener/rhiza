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
- serves change set `01M6FY9EXET6R8QN6GPEJTCF7B` (Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists)
- serves change set `01M6GDJW3R8PXW7QBAZVH0QS4W` (Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix')
- serves change set `01M6GJ9GBCQD9GFCDSNR4BB13Q` (Serialize the three bug change sets titled Fix, so their ship runs never share cs/fix and rhiza.fix: the test-scripts fix waits for the resume fix)
- serves change set `01M6H3TC786A6WVKRXC0WQ1N3K` (Agent panes close when their session finishes: every driver closes a finished session's pane at once (review checks, critic, fixers, plan steps), not at close-out)
- serves change set `01M6M01GYQHR9VMCZ54W5THARK` (Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped)
- serves change set `01M6MA6P7KKMQPFBM6ZD6EMDM3` (Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
