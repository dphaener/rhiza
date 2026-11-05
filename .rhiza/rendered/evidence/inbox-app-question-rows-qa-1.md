# inbox-app-question-rows-qa-1

evidence · `01M6PS0F53Y1A6TNSST09JEE3G` · actor [claude-code-4df78838](../actor/claude-code-4df78838.md)

## result

FAIL

## statement

QA at a2f3e32b with a binary built from the branch (go build -trimpath), in a scratch repo: rhiza tui app inbox in tmux at 60-160 columns: the list rows are one line each, truncated with an ellipsis ('Which disk should hold the … when the machin…'), not wrapped, and do not say what they block; the blocker row and answer form show 'hold the … when' because scrub replaced the substring 'store' (the blocked capability's handle) inside the question's prose. The gate detail panel (after answering) wrapped and showed no handle or ID except 'next: rhiza run resume …'.

## Links

### Out

- verifies [inbox-app-question-rows](../capability/rhiza/inbox-app-question-rows.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
