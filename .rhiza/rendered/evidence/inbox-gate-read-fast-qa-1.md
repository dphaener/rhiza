# inbox-gate-read-fast-qa-1

evidence · `01M6S3X8CD44KNTX9A966BYQ9K` · actor [claude-code-8eebcc9a](../actor/claude-code-8eebcc9a.md)

## result

PASS

## statement

QA at ed87c504, built binary (go build -trimpath), scratch clone of this repo (11,596 log entries, worktree index present), a draft plan change set with four open questions and a hand-written plan run journal stopped at its gate on them, worktree = the clone. rhiza inbox listed all four questions with full statements: new binary 1.1-1.7 s wall (load average 11-17 on 4 cores; 0.4-0.75 s with no gated run), merge-base binary 26 s with the same output. CPU profile: inboxGateRead 0.15 s, no store.ReadLog; the rest of its cost is writeIndex.changeSets decoding every folded change set (~0.76 s CPU over goroutines).

## Links

### Out

- verifies [inbox-gate-read-fast](../capability/rhiza/inbox-gate-read-fast.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
