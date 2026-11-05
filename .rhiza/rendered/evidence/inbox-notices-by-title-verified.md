# inbox-notices-by-title-verified

evidence · `01M6PRRE2FEEC8J4R2DAJNFEAA` · actor [claude-code-67f73741](../actor/claude-code-67f73741.md)

## result

PASS

## statement

At 1ba53138 I read the diff from 623c4061 and ran go test -trimpath ./internal/api/ -run 'TestInboxNotice|Daemon|Inbox' -count=1 (ok), go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden -count=1 (ok), plus go vet and gofmt -l on internal/api (clean). TestInboxNoticeNamesNoHandleOrID in internal/api/daemon_notice_test.go is a table over inboxNotice with a case for each of the 14 inbox kinds in inbox.go, and it fails if any kind has no case. Each item carries IDs and handles in every field (ID, Run, Change, Other, Handle, Blocks, Target, Session, TabBlock, Next), and the test asserts that no notice contains any of them. The overlap case expects the other change set's title, with 'another change set' when there is no title. The question case expects the question text. The open cases expect the first question's restatement, falling back to its statement. In production, OtherTitle is filled in daemonNotify and in runInbox. The open item's Questions come from inboxGateQuestions in daemonNotify. Not covered by this leaf: daemonNotify and daemonHook still prefix the notice with the change ID when no journal has a title for it. That prefix is set outside inboxNotice and predates this change.

## Links

### Out

- verifies [inbox-notices-by-title](../capability/rhiza/inbox-notices-by-title.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
