# inbox-app-question-rows-verified-2

evidence · `01M6PVSFSR02CBQMS79KC3V38G` · actor [claude-code-ad45ddcc](../actor/claude-code-ad45ddcc.md)

## result

PASS

## statement

Verifier at eeb31b73 (attempt 2). go test -trimpath -count=1 -v ./internal/tui/inboxapp/: all 34 tests pass, including TestBlockerShowsEveryWordWrapped (300+ char blocker, every word present once opened at width 70, no line wider than 70, 'Blocks: Web UI'), TestGateRowShowsQuestionsAndWhatTheyBlock (row shows the first question's restatement; opened view shows each question with 'Blocks: Pick storage' and 'Blocks: Scale out'), TestScrubKeepsWordsAndDropsOnlyWholeIDs (blocker on handle 'store' keeps 'hold the store when' in the question and 'Blocks: Kept store reachable'), TestNoSurfaceShowsAnIDOrHandle (every kind's row, form, detail, send-back and resolve form shows no handle, ULID or short ID), TestOverlapRowNamesBothChangeSetsAndTarget. go test ./internal/api -run 'Inbox|Overlap' ok; go test ./cmd/rhiza -run 'TestScripts/(describe-golden|inbox.*)$' ok. Read the diff: scrub now replaces only whole id-like tokens (digit, -, _ or .), overlap rows name the other change set and target statement, target_title added to InboxItem and the describe golden. A whole ./internal/api run hit the 600s go test timeout at load average ~10 (machine contention); the inbox-scoped subset passes in 15s.

## Links

### Out

- verifies [inbox-app-question-rows](../capability/rhiza/inbox-app-question-rows.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
