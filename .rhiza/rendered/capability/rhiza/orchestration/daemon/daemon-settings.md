# daemon-settings

capability · `01M6CS1WF6WPPJTPC8X1DHYGD8` · parent [daemon](../daemon.md)

## acceptance

daemon-settings.txtar: daemon settings lists every key with its default, scope and source default; daemon set rhiza.daemon.maxRuns 1 then settings shows 1 from rhiza.toml; an unknown key or bad value exits 2 writing nothing; with maxRuns set to 1 between two daemon --once polls, the second starts at most one run; a --max-runs flag wins over config and settings shows source flag in daemon status; describe-golden.txtar is regenerated

## statement

The daemon reads its settings every poll, the machine ones (paused, max runs, minimum free disk, Go cache trimming, permission mode) from the root's rhiza.toml and the project ones (plan, accept, ship and signoff autonomy) from each project's git config, as root-daemon-settings-scope splits them, a command-line flag winning for that run; rhiza daemon settings lists each with value, scope and source and rhiza daemon set validates and writes one to its scope

## Links

### Out

- decided-by [settings-in-git-config](../../../../decision/settings-in-git-config.md)

### In

- verifies [daemon-settings-verified](../../../../evidence/daemon-settings-verified.md)
- verifies [daemon-settings-qa-1](../../../../evidence/daemon-settings-qa-1.md)
- verifies [daemon-settings-verified-6](../../../../evidence/daemon-settings-verified-6.md)
- verifies [daemon-settings-qa-3](../../../../evidence/daemon-settings-qa-3.md)
- verifies [pr-187-live-daemon-settings](../../../../evidence/pr-187-live-daemon-settings.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
