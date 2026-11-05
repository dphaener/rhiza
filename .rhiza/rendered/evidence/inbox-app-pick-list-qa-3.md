# inbox-app-pick-list-qa-3

evidence · `01M6NF65RF2H49F56EJYEQ413W` · actor [claude-code-f6175daa](../actor/claude-code-f6175daa.md)

## result

PASS

## statement

QA at 851e0583 with the built binary (rhiza tui app inbox) in tmux on a real plan run stopped at its gate (plan-workflow-open fixture, fake rex/claude/gh), run from the main checkout: the open item's row shows kind, change set title, step and 'open questions: park-disk'; Enter opens the per-question answer form; typing 'Exports go on the ssd' and Enter wrote Decision answer-park-disk (quote and statement 'Exports go on the ssd', answers park-disk) into the planned change set in the run's worktree, queued the run in continues.json, and the app exited 0. Resuming the run then passed the gate and went on to open-pr; after a retry it opened one PR and the inbox said Nothing waits on you. Cycle 2's gate finding is fixed. go test -trimpath ./internal/tui/inboxapp/... and TestScripts plan-workflow-open, tui-inbox-app and question-answer pass.

## Links

### Out

- verifies [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
