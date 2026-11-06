# milestone-29

intent · `01M6RN3786W5DAVS69HMZ28RD9`

## milestone

29

## statement

Rhiza manages and builds several projects on one machine from one root directory: projects are registered under the root, one root daemon builds them all under one run cap and one usage budget, the CLI, web UI and TUI apps show every project and each on its own, and this VM runs rhiza and nebari that way.

## Links

### Out

- decided-by [m29-one-root-daemon](../decision/m29-one-root-daemon.md)
- decided-by [m29-plan-all-now](../decision/m29-plan-all-now.md)

### In

- serves change set `01M6RN3M8EC8FEDEJ8X6YVCPZ4` (Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project)
- serves change set `01M6RN3RCGAZ9RNTWTZPKFM61M` (Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service)
- serves change set `01M6RN3VGZ7D1JK4FRPNX21WHX` (Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one)
- serves change set `01M6RN3ZBGYDW6KPBSES9A5JXS` (Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher)
- serves change set `01M6RN441P2Z865ZVKAS7X7GMG` (Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox)
- serves change set `01M6RN474B5EXMX8NYHE7F62RV` (This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
