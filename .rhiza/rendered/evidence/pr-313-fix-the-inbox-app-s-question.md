# pr-313-fix-the-inbox-app-s-question

evidence · `01M6Q0W5D7WE7R4P6W6QSPK4PH` · actor [dphaener](../actor/dphaener.md)

## head

097c398effb0713b7bf9d0e03158e16f4614a41f

## pr

https://github.com/dphaener/rhiza/pull/313

## result

1 check passed: check SUCCESS

## statement

Pull request #313 (https://github.com/dphaener/rhiza/pull/313) passed CI at 097c398

## status

verified

## Links

### Out

- verifies [inbox-question-text](../capability/rhiza/inbox-question-text.md)
- verifies [owner-question-advisor](../capability/rhiza/owner-question-advisor.md)
- verifies [inbox-app-question-rows](../capability/rhiza/inbox-app-question-rows.md)
- verifies [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)
- verifies [inbox-notices-by-title](../capability/rhiza/inbox-notices-by-title.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
