# inbox-app-answer-shows-question-qa-1

evidence · `01M6PS0F39EPE5CB58K38DPND6` · actor [claude-code-4df78838](../actor/claude-code-4df78838.md)

## result

FAIL

## statement

QA at a2f3e32b with a binary built from the branch (go build -trimpath), in a scratch repo: rhiza tui app inbox driven in tmux. The gate form for the two-question plan run showed '1 of 2: <restatement>', Blocks, Recommended, Why and '2 of 2: <statement>' above 'Your answer to question N of 2', no handles; typing two answers and submitting wrote decisions answer-park-disk (quote 'Use the SSD') and answer-park-format (quote 'JSON') to the right questions (resume then failed only because the fake journal's workflow is empty). But the blocker's answer form at 60/80/90/160 columns showed the 300-char question on ONE line cut at the pane edge, not wrapped, until the pane was resized; and the word 'store' (the handle of what it blocks) was replaced by '…' in the question and in 'Blocks: Kept … reachable'.

## Links

### Out

- verifies [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
