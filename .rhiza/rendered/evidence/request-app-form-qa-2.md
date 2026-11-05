# request-app-form-qa-2

evidence · `01M6NBH96D3XRAR9WBFKR4R5SZ` · actor [claude-code-66d44bb0](../actor/claude-code-66d44bb0.md)

## result

PASS

## statement

QA cycle 2 at 1bd20b12 with a scratch build (go build -trimpath ./cmd/rhiza) in a fresh rhiza init repo, driving `rhiza tui app request` in a real tmux pty: the form drew a Kind Select (Feature, Bug, Milestone; Feature default) and a Request TextArea. Picking Milestone and pressing Enter on empty text showed '✗ A request needs its text' and queued nothing; typing 'VM milestone' and Enter printed the queued request (kind milestone) and exited 0. Picking Bug and submitting '--json   broken' queued kind bug, text '--json broken' (dash-leading text kept, whitespace collapsed), exit 0, listed by `rhiza request ls --json`. Typing 'dark mode' then Esc exited 0 and queued nothing. `tui app request --once` exits 2; the old `tui app requests` exits 4 'unknown app'. The slice shows requests-app and requests-app-milestone proposed for retirement in this change set. go test -trimpath ./internal/tui/requestapp/... and ./cmd/rhiza -run 'TestScripts/(tui-request-app|rex-install)$' pass.

## Links

### Out

- verifies [request-app-form](../capability/rhiza/request-app-form.md)

## History

- `01M6MYNDJ4HAJF028CEAE7SYQN` Request app: file a feature, bug or milestone request in one rexuikit Form, then close; rebuilt on the modal shell over the request operation, replacing internal/tui/requestsapp and retiring the capabilities of the old requests app it replaces
