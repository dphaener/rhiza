# m29-surfaces-across-projects

intent · `01M6RN379J7YGMNFSCC3MZ2RGV` · parent [milestone-29](../milestone-29.md)

## statement

The CLI, the web UI and the Rex TUI apps show every project at once, with a merged inbox, and each project on its own.

## Links

### Out

- decided-by [m29-web-all-projects](../../decision/m29-web-all-projects.md)
- decided-by [m29-cli-and-tui-now](../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-tui-current-project](../../decision/m29-tui-current-project.md)

### In

- serves change set `01M6RN3VGZ7D1JK4FRPNX21WHX` (Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one)
- serves change set `01M6RN3ZBGYDW6KPBSES9A5JXS` (Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher)
- serves change set `01M6RN441P2Z865ZVKAS7X7GMG` (Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
