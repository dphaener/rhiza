# inbox-app-question-rows-verified

evidence · `01M6PQG15E9MXZ1K83K1B43KCH` · actor [claude-code-5b09a2c5](../actor/claude-code-5b09a2c5.md)

## result

PASS

## statement

At dc8bf497 ran go test -trimpath -count=1 ./internal/tui/inboxapp/... (ok; go vet and gofmt clean, go build ./... ok). With the package's fake runner: TestBlockerShowsEveryWordWrapped opens a blocker whose question is over 300 characters at width 70 and finds every word in the view, no line wider than 70, plus 'Blocks: Web UI'; TestGateRowShowsQuestionsAndWhatTheyBlock checks the gate row shows the restatement (not the statement) and the opened gate shows each question's text with 'Blocks: Pick storage' / 'Blocks: Scale out', recommendation and reasons; TestNoSurfaceShowsAnIDOrHandle renders the rows, the opened form, dialog or detail panel of every item kind (question, blocker, gate with and without questions, signoff, failed, stalled, overlap, accept) plus the send-back and resolve forms and finds none of the item, run, session, change set, node IDs or handles. Read the diff: run IDs dropped from the signoff dialog and abandon confirm, gate questions taken from the structured questions field, agent text scrubbed of the item's IDs, list rows stay one line (full wrapped text shows once the item is opened).

## Links

### Out

- verifies [inbox-app-question-rows](../capability/rhiza/inbox-app-question-rows.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
