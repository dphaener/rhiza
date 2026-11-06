# inbox-gate-read-fast-qa-2

evidence · `01M6S5M9NECEX5C0HQR6452M6B` · actor [claude-code-9d02afc1](../actor/claude-code-9d02afc1.md)

## result

PASS

## statement

QA at 9efc1e89 with the built binary in the same scratch clone (11,630 log entries, load average 14 on 4 cores): rhiza inbox with a plan run stopped at its gate on four open questions listed all four with full statements and next: rhiza run resume in 2.1 s wall (1.0 s user CPU), versus ~26-36 s for the log-decoding reads before this change set.

## Links

### Out

- verifies [inbox-gate-read-fast](../capability/rhiza/inbox-gate-read-fast.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
