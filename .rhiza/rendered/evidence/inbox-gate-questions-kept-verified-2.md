# inbox-gate-questions-kept-verified-2

evidence · `01M6S489G2HYBEP6NYVHNRJ3WA` · actor [claude-code-a725333b](../actor/claude-code-a725333b.md)

## result

PASS

## statement

Verifier at 968be804 (attempt 2): `go test -trimpath ./internal/api -run TestInboxGateQuestions -count=1 -v` passes TestInboxGateQuestionsKept (case a: worktree with no .rhiza, stop detail 'kept-a, kept-b' -> both handles listed in order, statements empty) and TestInboxGateQuestionsUnresolved (case b: readable worktree, a proposed change set serving the plan's intent holds question kept-resolved; stop detail 'kept-resolved, kept-nowhere' -> both handles listed, kept-resolved with 'Statement of kept-resolved', kept-nowhere empty), plus TestInboxGateQuestionsBudget. `go test -trimpath ./cmd/rhiza -run TestScripts/inbox-gate-questions -count=1` passes. Read the diff: test-only, and inboxGateQuestions in internal/api/inbox.go seeds one entry per stop-detail handle before the worktree read, so either fallback keeps every handle.

## Links

### Out

- verifies [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
