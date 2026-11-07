# tui-merged-inbox

capability · `01M6RN86Q8DXQ5T5PCRHDSBRMM` · parent [rhiza](../rhiza.md)

## acceptance

Inbox app tests with two projects' items check the project tag on each row and that answering an item from project B writes into B, opened from project A.

## statement

The Inbox app lists every project's inbox items, each row tagged with its project, and answering, approving, retrying or abandoning an item acts in that item's project.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-tui-current-project](../../decision/m29-tui-current-project.md)

### In

- verifies [tui-merged-inbox-verified](../../evidence/tui-merged-inbox-verified.md)
- verifies [tui-merged-inbox-qa-1](../../evidence/tui-merged-inbox-qa-1.md)
- verifies [pr-340-project-aware-tui-apps](../../evidence/pr-340-project-aware-tui-apps.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
