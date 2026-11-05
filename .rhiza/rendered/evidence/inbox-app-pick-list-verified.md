# inbox-app-pick-list-verified

evidence · `01M6NAEB62N57BSJ3RG0JHY33M` · actor [claude-code-19373354](../actor/claude-code-19373354.md)

## result

PASS

## statement

Verifier claude-code-19373354 at 14bb597f: read the diff of internal/tui/inboxapp (one rexuikit ui.List with a row per item in the inbox operation's order; ActionFor maps question/open/blocker to answer, signoff to signoff, failed/gave-up/skipped to resolve, accept to accept, everything else to show; show opens an InfoPanel with the detail and calls inbox.open on the item's run). Ran go test -trimpath -count=1 ./internal/tui/inboxapp/... (ok): TestRowPerItemInOrder (six rows in order with kind, change set title and what it waits for), TestEnterOpensActionPerKind (Enter on question, signoff, failed and accept items calls the Open hook with that action and no runner op), TestActionFor (all 13 kinds incl. stalled, overlap, disk-low, usage-reserve, main-stale -> show), TestOtherKindShowsDetailAndOpensTab (Enter on stalled shows its detail and calls inbox.open h8), TestEmptyInbox (one empty-state line). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(tui-inbox-app|inbox-open|tui-.*)$' (ok): tui-inbox-app cmps rhiza tui app inbox --once with rhiza inbox --json byte for byte for an empty inbox and a failed run. go build ./... and go vet on internal/tui and cmd/rhiza clean.

## Links

### Out

- verifies [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
