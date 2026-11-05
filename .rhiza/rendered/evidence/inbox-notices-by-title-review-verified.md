# inbox-notices-by-title-review-verified

evidence · `01M6PVNKTKR4WBDYRQK637B5XF` · actor [claude-code-af44f6f9](../actor/claude-code-af44f6f9.md)

## result

PASS

## statement

At 791b7139 (review cycle 2 fix for model finding f3): read the diff. daemonNotify now fills titles for change sets with no run from the fold (DisplayTitle), and noticeSubject falls back to 'a change set', never the ID, in both the Rex notice and daemonHook. Ran go test -trimpath -count=1 ./internal/api/ -run 'TestInboxNotice|TestDaemonNotify|TestInbox|Daemon' (ok): TestInboxNoticeNamesNoHandleOrID has one case per inbox kind with IDs and handles in every field, none appear, overlap holds the other change set's title, question and open hold the question's text; TestDaemonNotifyOverlapByTitle shows an overlap between two run-less change sets notified as 'Rename N: an overlap with Rename M' with neither ID. Ran go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/daemon-notify (ok), whose hook expects 'a change set' in place of the ID.

## Links

### Out

- verifies [inbox-notices-by-title](../capability/rhiza/inbox-notices-by-title.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
