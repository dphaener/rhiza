# request-app-form-qa-1

evidence · `01M6NAT9BTBXXVHF0XW5X7A65E` · actor [claude-code-bb6de9d1](../actor/claude-code-bb6de9d1.md)

## result

PASS

## statement

QA at 1021835c with a scratch build (go build -trimpath ./cmd/rhiza) in a fresh rhiza init repo, driving `rhiza tui app request` in a real tmux pty: the form drew a Kind Select (Feature, Bug, Milestone; Feature default) and a Request TextArea. Picking Bug and pressing Enter on empty text showed '✗ A request needs its text' and queued nothing. Typing '--json output   is wrong' and Enter printed the queued request JSON (kind bug, text '--json output is wrong', whitespace collapsed, dash-leading text kept) and exited 0; `rhiza request ls --json` listed it queued. Picking Milestone and submitting 'VM milestone' queued kind milestone. Typing 'dark mode' then Esc exited 0 and queued nothing. `tui app request --once` exits 2 with 'has no --once'; the old name `tui app requests` exits 4 'unknown app'. go test -trimpath ./internal/tui/requestapp/... and ./cmd/rhiza -run 'TestScripts/(tui-request-app|rex-install)$|TestRequestsApp' pass.

## Links

### Out

- verifies [request-app-form](../capability/rhiza/request-app-form.md)

## History

- `01M6MYNDJ4HAJF028CEAE7SYQN` Request app: file a feature, bug or milestone request in one rexuikit Form, then close; rebuilt on the modal shell over the request operation, replacing internal/tui/requestsapp and retiring the capabilities of the old requests app it replaces
