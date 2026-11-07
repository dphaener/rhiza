# tui-merged-inbox-verified

evidence · `01M6SPX8SG6SK1V0X8D4MMQEQ1` · actor [claude-code-d5cac357](../actor/claude-code-d5cac357.md)

## result

PASS

## statement

At 157661bf, read the diff: inboxapp.NewMergedModel loads every registered project's inbox through a ProjectRunner, sorts items by when they began to wait, tags each row '<project> · ...', and every write (question.answer, session.message, accept, run.approve, run.send-back, run.resume --retry, run.abandon, inbox.open) goes through act(), which runs in the picked item's project; cmd/rhiza appSwitcher uses it when projectctx.Folders finds registered projects, chdir'ing into the item's project folder per call (newCtx resolves the root from the working directory). Ran go test -trimpath ./internal/tui/inboxapp/... -run Merged -v: TestMergedRowsTaggedWithProject (two projects alpha and beta, rows 'beta · failed · Beta work' then 'alpha · question · Alpha work'), TestMergedAnswerWritesIntoItemsProject (beta's failed item retried with the last call run.resume --retry RB in project beta, not alpha) and TestMergedQuestionAnswerGoesToItsProject (alpha's question answered via session.message in alpha) all PASS. Also go test -trimpath ./internal/tui/..., go vet ./cmd/rhiza ./internal/tui/..., and go test -trimpath ./cmd/rhiza -run 'TestScripts/(tui-|inbox-)' plus the TUI unit tests: all ok.

## Links

### Out

- verifies [tui-merged-inbox](../capability/rhiza/tui-merged-inbox.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
