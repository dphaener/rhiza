# inbox-app-question-rows

capability · `01M6PN1T0KHDFWC5KNHJJHSQZT` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: a blocker with a 300-character statement renders every word across wrapped lines at the modal's width once opened; a gate row shows its first question's text, and opening it shows each question's text and what each blocks; a blocker on a node whose handle is an ordinary word keeps that word in its question and Blocks line; a test renders every item kind's row, form and detail panel and asserts no handle, ULID or short ID from the items appears.

## statement

In the Inbox app a question, blocker or gate row shows its kind, its change set title and the question's restatement (else its full statement) on one line, cut to fit; opening the item shows the question wrapped over as many lines as it needs and what it blocks by title. No row, form or detail panel shows a handle or an ID (item, run, change set or node), and a handle that is an ordinary word stays in the text.

## Links

### Out

- decided-by [owner-surfaces-full-text](../../decision/owner-surfaces-full-text.md)

### In

- verifies [inbox-app-question-rows-verified](../../evidence/inbox-app-question-rows-verified.md)
- verifies [inbox-app-question-rows-qa-1](../../evidence/inbox-app-question-rows-qa-1.md)
- verifies [inbox-app-question-rows-verified-2](../../evidence/inbox-app-question-rows-verified-2.md)
- verifies [inbox-app-question-rows-qa-2](../../evidence/inbox-app-question-rows-qa-2.md)
- verifies [pr-313-fix-the-inbox-app-s-question](../../evidence/pr-313-fix-the-inbox-app-s-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
