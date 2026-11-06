# inbox-gate-questions-kept-verified

evidence · `01M6S2NGYHZF19GEFQJAFA04C7` · actor [claude-code-6386b61f](../actor/claude-code-6386b61f.md)

## result

PASS

## statement

Verifier claude-code at 6d9163c7: read the diff (inboxGateQuestions now seeds one question per stop-detail handle, statement empty, and only replaces them when inboxGateRead reads the worktree; advice matching moved after, by handle or resolved ID). Ran go test -trimpath ./internal/api -run TestInboxGateQuestions -count=1 -v: TestInboxGateQuestionsKept (worktree with no .rhiza, stop on kept-a, kept-b -> both handles, statements empty) PASS, TestInboxGateQuestionsBudget PASS. Ran go test -trimpath ./cmd/rhiza -run TestScripts/inbox-gate-questions -count=1: PASS. Also checked end to end with a scratch build in a hermetic repo: a plan run stopped on 'park-disk, nowhere-q' (nowhere-q resolves in no planned change set) shows the open item with 'question: Which-disk (park-disk)' and 'question:  (nowhere-q)'; with the run worktree pointed at /nonexistent, inbox --json lists both handles with statement empty. Note: the Go test covers the unreadable worktree via inboxGateQuestions (the function that fills open items, inbox.go:634, daemon.go:403); I checked the unresolved-handle case by hand only, because no Go test covers it.

## Links

### Out

- verifies [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
