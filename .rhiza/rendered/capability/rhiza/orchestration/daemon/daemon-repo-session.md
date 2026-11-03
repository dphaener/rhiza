# daemon-repo-session

capability · `01M6FES3PKQCSNG9SEECK4G1R7` · parent [daemon](../daemon.md)

## acceptance

A testscript runs rhiza daemon --once against a fake Rex with REX_SESSION unset: when rex session ls lacks rhiza-<repo> the daemon creates it once and every detached run driver it starts carries REX_SESSION=rhiza-<repo>; with the session present it creates none; with REX_SESSION set to another session those drivers still carry rhiza-<repo>

## statement

Every Rex tab the daemon's runs open (agent tabs; plan-request runs and run drivers start headless, outside Rex, by daemon-headless-run-drivers) opens in the Rex session rhiza-<repo>, named after the main checkout's folder, which the daemon creates on its Rex server when missing and passes as REX_SESSION to every run driver it starts, so it needs no Rex tab or REX_SESSION of its own

## Links

### In

- verifies [daemon-repo-session-verified](../../../../evidence/daemon-repo-session-verified.md)
- verifies [daemon-repo-session-qa-1](../../../../evidence/daemon-repo-session-qa-1.md)
- verifies [pr-200-the-daemon-runs-as-a-systemd](../../../../evidence/pr-200-the-daemon-runs-as-a-systemd.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
