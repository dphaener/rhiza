# inbox-question-text

capability · `01M6PN1SXMJV5FCGDM5YXJ7RQX` · parent [rhiza](../rhiza.md)

## acceptance

A new testscript inbox-gate-questions (plus the existing inbox scripts and describe-golden, regenerated) on a plan run stopped at its gate with two open questions and a blocker: rhiza inbox --json lists each gate question with its handle, full statement and blocks title, and the blocker with the title of what it blocks; rhiza inbox prints each full statement and, after a recorded recommendation, the restatement and recommendation.

## statement

The inbox gives every owner question its full text: a gate (open) item lists each open question it stopped on with its full statement, the advisor's plain restatement and recommendation once written, and what it blocks by statement or change set title; a blocker item names what it blocks the same way; --json keeps each question's handle for the commands, and rhiza inbox's text prints each question's full statement and recommendation, not its handle alone.

## Links

### Out

- decided-by [owner-surfaces-full-text](../../decision/owner-surfaces-full-text.md)
- decided-by [owner-question-wording-by-advisor](../../decision/owner-question-wording-by-advisor.md)

### In

- verifies [inbox-question-text-verified](../../evidence/inbox-question-text-verified.md)
- verifies [inbox-question-text-qa-1](../../evidence/inbox-question-text-qa-1.md)
- verifies [inbox-question-text-qa-2](../../evidence/inbox-question-text-qa-2.md)
- verifies [pr-313-fix-the-inbox-app-s-question](../../evidence/pr-313-fix-the-inbox-app-s-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
