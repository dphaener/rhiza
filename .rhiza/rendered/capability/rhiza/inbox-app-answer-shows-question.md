# inbox-app-answer-shows-question

capability · `01M6PN1T22SZY3QGXYD783Z3JE` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: a question's answer form shows its full statement and recommendation above the TextArea; a gate form with two questions shows each one's statement above its own TextArea with no handle in the view, and submitting calls question.answer with each question's handle and its typed text.

## statement

The Inbox app's answer form, and its gate form for each open question, shows the full question wrapped above its text box (the advisor's restatement and recommendation once written, else the statement), never a handle as the label, while the answer still goes to the right question.

## Links

### Out

- decided-by [owner-surfaces-full-text](../../decision/owner-surfaces-full-text.md)
- decided-by [owner-question-wording-by-advisor](../../decision/owner-question-wording-by-advisor.md)

### In

- verifies [inbox-app-answer-shows-question-verified](../../evidence/inbox-app-answer-shows-question-verified.md)
- verifies [inbox-app-answer-shows-question-qa-1](../../evidence/inbox-app-answer-shows-question-qa-1.md)
- verifies [inbox-app-answer-shows-question-verified-2](../../evidence/inbox-app-answer-shows-question-verified-2.md)
- verifies [inbox-app-answer-shows-question-qa-2](../../evidence/inbox-app-answer-shows-question-qa-2.md)
- verifies [pr-313-fix-the-inbox-app-s-question](../../evidence/pr-313-fix-the-inbox-app-s-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
