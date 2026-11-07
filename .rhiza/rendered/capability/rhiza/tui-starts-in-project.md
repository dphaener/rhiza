# tui-starts-in-project

capability · `01M6RN86K392PSD63NH20X4D85` · parent [rhiza](../rhiza.md)

## acceptance

TUI tests submit the Request app in project B's folder and find the request queued in B; at the root the app shows the picker, Esc starts nothing, and picking A queues it in A.

## statement

Everything a TUI app starts (a milestone, a change set, a plan or a request) is started in the app's current project, and opened with no current project the app asks for one from a picker first, starting nothing until one is picked.

## Links

### Out

- decided-by [m29-tui-current-project](../../decision/m29-tui-current-project.md)
- decided-by [m29-tui-project-from-pane](../../decision/m29-tui-project-from-pane.md)

### In

- verifies [tui-starts-in-project-verified](../../evidence/tui-starts-in-project-verified.md)
- verifies [tui-starts-in-project-qa-1](../../evidence/tui-starts-in-project-qa-1.md)
- verifies [pr-340-project-aware-tui-apps](../../evidence/pr-340-project-aware-tui-apps.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
