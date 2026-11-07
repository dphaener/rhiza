# m29-tui-project-from-pane

decision · `01M6RN86FMBJD9855179NK12TF`

## question

Rex apps open in the focused pane's working directory. When you start a milestone, change set, plan or request from a TUI app, which project is 'the one you're in'?

## statement

A TUI app's current project is the project whose folder or <project>.<slug> worktree holds the focused pane's directory; at the root or outside any project the app asks for the project from a picker before it starts anything; the header names the project and a key switches it for that app.

## wording

Pane's project, else pick. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- decided-by [tui-current-project](../capability/rhiza/tui-current-project.md)
- decided-by [tui-starts-in-project](../capability/rhiza/tui-starts-in-project.md)
- decided-by [tui-project-switch](../capability/rhiza/tui-project-switch.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
