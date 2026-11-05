# inbox-app-pick-list-qa-2

evidence · `01M6NDHH0R41172WECTTGPPMRC` · actor [claude-code-3293f9c2](../actor/claude-code-3293f9c2.md)

## result

FAIL

## statement

QA at 433db495 with the built binary in tmux on testscript scratch repos: one List row per item in rhiza inbox order with kind, change set title, step and what it waits for (failed, question, signoff, blocker, disk-low; skipped and accept; open); rhiza tui app inbox --once is byte-identical (cmp) to rhiza inbox --json on the accept and open repos; an empty inbox shows 'Nothing waits on you'. Enter opens the action on question, blocker, signoff, failed, skipped and accept (Confirm, default No; Yes accepted the change set and it left the inbox); Enter on disk-low shows its detail view. But on an open (gate) item from a real stopped plan run, the per-question form submits question.answer against the main repo and the run's record change set and always fails ('answer failed: no question "park-disk"'): the parked question lives only in the plan run's worktree and planned change set, so a gate can never be answered from the app.

## Links

### Out

- verifies [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
