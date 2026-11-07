# tui-project-switch

capability · `01M6RN86NDXBMPTCAW1V6MX121` · parent [rhiza](../rhiza.md)

## acceptance

TUI tests open an app in project A, switch to B through the picker, and see the header and the listed data change to B's.

## statement

In every TUI app one key opens a project picker of the registered projects that switches the app's current project, and the app reloads its view for the picked project.

## Links

### Out

- decided-by [m29-tui-project-from-pane](../../decision/m29-tui-project-from-pane.md)
- decided-by [m29-cli-and-tui-now](../../decision/m29-cli-and-tui-now.md)

### In

- verifies [tui-project-switch-verified](../../evidence/tui-project-switch-verified.md)
- verifies [tui-project-switch-qa-1](../../evidence/tui-project-switch-qa-1.md)
- verifies [pr-340-project-aware-tui-apps](../../evidence/pr-340-project-aware-tui-apps.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
