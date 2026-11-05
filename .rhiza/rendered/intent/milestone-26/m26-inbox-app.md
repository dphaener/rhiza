# m26-inbox-app

intent · `01M6MZ0BVXSHHDCB4GXVABYTKE` · parent [milestone-26](../milestone-26.md)

## statement

From the Inbox app the owner picks an inbox item and does its one action (answers a question or gate, signs off a run, retries or abandons a failed run), and the app closes.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-run-actions](../../decision/tui-run-actions.md)

### In

- serves change set `01M6MYN9GRWHRJSHA31D9S838H` (Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces)
- serves change set `01M6PMSN1PKES7GJYAEXNZSN13` (Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.)

## History

- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
