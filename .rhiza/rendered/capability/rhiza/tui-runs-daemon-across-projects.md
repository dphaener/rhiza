# tui-runs-daemon-across-projects

capability · `01M6RN86SPHVRJW809SPKT0ASP` · parent [rhiza](../rhiza.md)

## acceptance

TUI tests with two projects check the Runs app's project tags, and that toggling an autonomy stage in the Autonomy app writes the current project's git config and setting max runs writes rhiza.toml.

## statement

The Runs app shows every project's live runs tagged by project, and the Daemon and Autonomy apps show the machine settings once and the current project's settings, writing each to its scope.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-settings-split](../../decision/m29-settings-split.md)

### In

- verifies [tui-runs-daemon-across-projects-verified](../../evidence/tui-runs-daemon-across-projects-verified.md)
- verifies [tui-runs-daemon-across-projects-qa-1](../../evidence/tui-runs-daemon-across-projects-qa-1.md)
- verifies [pr-340-project-aware-tui-apps](../../evidence/pr-340-project-aware-tui-apps.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
