# owner-surfaces-full-text

decision · `01M6PN1SV5R599HFBNJ0N4C1BJ`

## quote

The full text would have been more helpful than the internal handle. No need to ever show the internal handle, it's nonsense to a user.

## statement

No owner-facing surface shows a node's internal handle or ID: show the question's full text and refer to things by their titles

## Links

### Out

- derived-from [inbox-question-fix-plan](../source/inbox-question-fix-plan.md)

### In

- decided-by [inbox-question-text](../capability/rhiza/inbox-question-text.md)
- decided-by [inbox-app-question-rows](../capability/rhiza/inbox-app-question-rows.md)
- decided-by [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)
- decided-by [inbox-notices-by-title](../capability/rhiza/inbox-notices-by-title.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
