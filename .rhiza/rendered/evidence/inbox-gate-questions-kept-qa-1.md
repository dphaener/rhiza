# inbox-gate-questions-kept-qa-1

evidence · `01M6S3X8H632BMG26DT3C8BY1T` · actor [claude-code-8eebcc9a](../actor/claude-code-8eebcc9a.md)

## result

PASS

## statement

QA at ed87c504 with the built binary in the same scratch clone: the gated plan run's journal pointed at worktree /nonexistent-qa; rhiza inbox still lists the open item with one 'question:  (qa-qN)' line per handle, statements empty, and next: rhiza run resume. The merge-base binary showed no question lines for that case. In the Inbox app (rhiza tui app inbox under tmux, 90x30) that item opened the gate answer form.

## Links

### Out

- verifies [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
