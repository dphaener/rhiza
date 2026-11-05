# inbox-question-text-verified

evidence · `01M6PQ12SR429CZK5WX752AZ6F` · actor [claude-code-6038ba9d](../actor/claude-code-6038ba9d.md)

## result

PASS

## statement

Verifier at fcc25e3e: read the diff (internal/api/inbox.go adds InboxItem.Questions []InboxGateQuestion {handle, statement, blocks, recommendation, reasons} read from the run's worktree for open items, and BlocksTitle on blockers, printed by Human). Ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(inbox|describe-golden)' -v: all 13 inbox scripts, including the new inbox-gate-questions, and describe-golden PASS. inbox-gate-questions sets up a plan run stopped at its gate on two open questions (park-disk, park-format) plus an accepted change set's blocker. rhiza inbox prints each full statement with its handle and what it blocks by statement, and the blocker as 'blocks: Kept (store)'. --json carries each question's handle and statement, blocks 'Exports-go-to-the-store' and blocks_title 'Kept'. After a journalled advisor recommendation, both text and --json show the recommendation and reasons. Also ran go test -trimpath ./internal/api/... (ok); gofmt -l and go vet on internal/api are clean. Note: a separate restatement field is not part of this leaf; the sibling advisor leaf adds it, and this leaf shows the advisor's recorded recommendation text and reasons.

## Links

### Out

- verifies [inbox-question-text](../capability/rhiza/inbox-question-text.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
