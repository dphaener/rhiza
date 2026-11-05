# owner-question-advisor-qa-1

evidence · `01M6PS0F12XAG7RGAAV8V7T2GH` · actor [claude-code-4df78838](../actor/claude-code-4df78838.md)

## result

PASS

## statement

QA at a2f3e32b with a binary built from the branch (go build -trimpath), in a scratch repo: the inbox-question-advisor workflow (command step adds question gate-q and journals a session-less asked entry, then stops) run via rhiza run start --foreground with fake claude/rex on PATH. The run logged 'advisor session ... on question gate-q', the fake advisor got the template text in its prompt and ran rhiza run recommend; the journal got one advisor and one recommended entry with asker = the question node ID and restated text; rhiza inbox then showed 'in plain words: Where should exports be stored?', the recommendation and reasons, and the statement stayed the template text.

## Links

### Out

- verifies [owner-question-advisor](../capability/rhiza/owner-question-advisor.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
