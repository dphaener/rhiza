# inbox-question-text-qa-1

evidence · `01M6PS0EZ6T5413EMQAGPH64NZ` · actor [claude-code-4df78838](../actor/claude-code-4df78838.md)

## result

PASS

## statement

QA at a2f3e32b with a binary built from the branch (go build -trimpath), in a scratch repo: an accepted change set 'Blocked' with question store-q (300-char statement) blocking capability store, and a plan change set 'Plan reports' with a run journal stopped at its gate on park-disk and park-format plus a journalled advisor recommendation on park-disk. rhiza inbox printed both gate questions' full statements with their handles, 'in plain words', 'blocks: Exports go to the store', the recommendation and reasons on park-disk, and the blocker's full 300-char statement with 'blocks: Kept store reachable (store)'. rhiza inbox --json gave questions[] with handle, statement, blocks, restatement, recommendation, reasons, and the blocker's blocks_title.

## Links

### Out

- verifies [inbox-question-text](../capability/rhiza/inbox-question-text.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
