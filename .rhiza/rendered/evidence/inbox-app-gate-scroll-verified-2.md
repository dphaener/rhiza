# inbox-app-gate-scroll-verified-2

evidence · `01M6S56G9Y3DAN3T5N68SFGV0S` · actor [claude-code-867d575d](../actor/claude-code-867d575d.md)

## result

PASS

## statement

Verifier at 3c0b06a0 (fixer attempt 2, qa f1). Read the diff: gateAnswerCtx now opens the run worktree's own index (index.DefaultPath + openWriteIndex, memoised models) for the gate's reads and keeps c.IndexPath for the answer's write, falling back to the log when no index opens, as inboxGateRead does. go test -trimpath ./internal/tui/inboxapp/... passed, incl. TestGateFormScrollsLongQuestions, TestGateUnreadableQuestionShowsPositionNotHandle, TestGateAnswersEachThenResumes, TestGateBlankAnswerSendsNothing, TestGateAnswerFailureReopensOnRemaining, TestGateResumeFailureDetailScrolls, TestDetailPanelScrolls. go test -trimpath ./internal/api -run 'Gate' passed (TestGateAnswerUsesWorktreeIndex, TestInboxGateQuestions{Budget,Kept,Unresolved}). In a scratch worktree with question.go reverted to 3c0b06a0~1, TestGateAnswerUsesWorktreeIndex fails ('decoded the full log 16 times'), so the test guards the fix.

## Links

### Out

- verifies [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
