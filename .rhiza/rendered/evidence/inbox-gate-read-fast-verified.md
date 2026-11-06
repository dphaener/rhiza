# inbox-gate-read-fast-verified

evidence · `01M6S31QDYHZ52GP9GKTHZA15H` · actor [claude-code-d0c13c24](../actor/claude-code-d0c13c24.md)

## result

PASS

## statement

At c30d3cb2 ran `go test -trimpath -count=1 ./internal/api -run Inbox -v`: all seven inbox tests PASS, including TestInboxGateQuestionsBudget (generated 3000-node log with an index; plan run at its gate with six open questions over three planned change sets). It asserts every question's statement, zero full-log decodes (countingCtx), each change set's model built at most once and every planned one (own plus three) exactly once via the builtModel hook in writeIndex.buildModel, no index catch-up decode after the read (LastCatchUp Decoded 0, not rebuilt), and the 1 s budget. Read the diff: inboxGateRead opens the worktree's index (index.DefaultPath) with an ix.models memo plus a local per-call modelOf memo. Mutation check in a scratch copy with both memos removed: the budget test FAILS with planned change sets folded 2-7 times, so the attempt 1 finding (no fold-count assertion) is fixed.

## Links

### Out

- verifies [inbox-gate-read-fast](../capability/rhiza/inbox-gate-read-fast.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
