# inbox-notices-by-title

capability · `01M6PN1T3JF18NJRVC33BQH6ZM` · parent [rhiza](../rhiza.md)

## acceptance

A table test in internal/api over inboxNotice, one case per inbox kind with items carrying IDs and handles: no notice contains any of them, the overlap notice holds the other change set's title, and the question and open notices hold the question's text.

## statement

The daemon's inbox notifications name change sets by title and questions by their restatement or full statement, never by handle or ID: an overlap names the other change set's title, and a gate stopped on open questions names the first question's text.

## Links

### Out

- decided-by [owner-surfaces-full-text](../../decision/owner-surfaces-full-text.md)

### In

- verifies [inbox-notices-by-title-verified](../../evidence/inbox-notices-by-title-verified.md)
- verifies [inbox-notices-by-title-review-verified](../../evidence/inbox-notices-by-title-review-verified.md)
- verifies [inbox-notices-by-title-qa-2](../../evidence/inbox-notices-by-title-qa-2.md)
- verifies [pr-313-fix-the-inbox-app-s-question](../../evidence/pr-313-fix-the-inbox-app-s-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
