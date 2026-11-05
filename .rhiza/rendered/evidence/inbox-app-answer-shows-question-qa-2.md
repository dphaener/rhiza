# inbox-app-answer-shows-question-qa-2

evidence · `01M6PW2NJ7GJCDD5APCPNEEXKC` · actor [claude-code-12e7dc02](../actor/claude-code-12e7dc02.md)

## result

PASS

## statement

QA cycle 2 at 61f23305 with a binary built from the branch (go build -trimpath), in a scratch repo with an accepted change set 'Keep exports in the store' blocked by a 300-char question on capability 'store' and by a leaf-failed template question with journalled advisor advice, plus a plan run stopped at its gate on two questions (one with advice). rhiza tui app inbox driven in tmux: the blocker's answer form wrapped the full question to the pane at 60, 80 and 160 columns with 'Blocks: Kept store reachable' (the word 'store' kept) above 'Your answer'; the leaf-failed blocker's form showed the restatement, 'Recommended: Keep the old path' and 'Why: Callers still use it' above the text box; the gate form showed '1 of 2: <restatement>', Blocks, Recommended, Why, '2 of 2: <statement>' above 'Your answer to question N of 2', no handles. Submitting wrote answer-park-disk ('Use the SSD') and answer-park-format ('JSON') to the right questions, and the blocker's answer wrote answer-store-q answering store-q in its change set.

## Links

### Out

- verifies [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
