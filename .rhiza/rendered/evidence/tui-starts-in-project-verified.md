# tui-starts-in-project-verified

evidence · `01M6SP965ABYTGB5BHA7ZHYDA5` · actor [claude-code-9bcf62bd](../actor/claude-code-9bcf62bd.md)

## result

PASS

## statement

At ce0597ee: go test -trimpath ./internal/tui/requestapp/... ./internal/tui/projectctx/... passes, including TestRequestQueuedInCurrentProject (app in project B submits; the runner is called for B), TestPickerAtRootThenQueuesInPicked (no project: Project picker shown, picking A queues nothing until the form submits, then queues in A) and TestPickerEscStartsNothing. End to end, a scratch build of rhiza driven in a pty over a scratch root with projects A and B: `rhiza tui app request` opened in B's folder queued the request in B/.git/rhiza/requests.json only; opened at the root it showed the Project picker, Esc closed it with nothing queued in A or B, and Enter (A) then a request queued it in A only. go vet and gofmt are clean. The diff wires cmd/rhiza's request app to projectctx.Current/Folders and chdirs to the picked project's folder before the in-process call, which resolves its repo from the working directory.

## Links

### Out

- verifies [tui-starts-in-project](../capability/rhiza/tui-starts-in-project.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
