# inbox-notices-by-title-qa-2

evidence · `01M6PW2NV3M2RANPN1QXQV93M7` · actor [claude-code-12e7dc02](../actor/claude-code-12e7dc02.md)

## result

PASS

## statement

QA cycle 2 at 61f23305 with a binary built from the branch (go build -trimpath), in a scratch repo with an accepted change set 'Keep exports in the store' blocked by a 300-char question on capability 'store' and by a leaf-failed template question with journalled advisor advice, plus a plan run stopped at its gate on two questions (one with advice). rhiza daemon --once --max-runs 0 with a fake notify hook (fake claude/rex/gh on PATH) gave 'Plan reports: a plan run stopped for open questions, the first: Choose where exports are stored' and 'a change set: a run failed' (that run's change set is not in the model): no handle or ID. go test ./internal/api -run 'Inbox|Notice' ok.

## Links

### Out

- verifies [inbox-notices-by-title](../capability/rhiza/inbox-notices-by-title.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
