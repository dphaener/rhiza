# root-daemon-settings-scope

capability · `01M6RN85C17S5VE1A6SNZAXK89` · parent [daemon](../daemon.md)

## acceptance

A testscript sets a machine key and a project key, finds the first in rhiza.toml and the second in that repo's git config, sees both in daemon settings --json with their scope and source, sees rhiza.daemon.trimGoCache listed with machine scope, and gets exit 2 for a project key set at the root with no --project; a daemon test with two projects, each with its own rhiza.eventCommand, runs each project's command only for that project's moves.

## statement

rhiza daemon settings lists every setting with its scope, machine (from rhiza.toml: maxRuns, usage budget and reserve, minFree, trimGoCache, web host and port, selfUpdate, permissionMode, paused) or project (from that repo's git config: autonomy per stage, project paused, per-project maxRuns, deploy check), and rhiza daemon set writes each key to its scope, taking --project or the current project for project keys. The hook commands rhiza.eventCommand and rhiza.notifyCommand are project settings: the root daemon runs a project's event command for that project's lifecycle moves and its notify command for that project's inbox, each read from that project's git config.

## Links

### Out

- decided-by [m29-settings-split](../../../../decision/m29-settings-split.md)

### In

- verifies [root-daemon-settings-scope-verified](../../../../evidence/root-daemon-settings-scope-verified.md)
- verifies [root-daemon-settings-scope-qa-1](../../../../evidence/root-daemon-settings-scope-qa-1.md)
- verifies [root-daemon-settings-scope-verified-5](../../../../evidence/root-daemon-settings-scope-verified-5.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
