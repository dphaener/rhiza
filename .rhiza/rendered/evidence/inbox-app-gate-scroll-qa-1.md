# inbox-app-gate-scroll-qa-1

evidence · `01M6S3X8PM1W21EGTCPF3K99A1` · actor [claude-code-8eebcc9a](../actor/claude-code-8eebcc9a.md)

## result

PASS

## statement

QA at ed87c504: rhiza tui app inbox (built binary) under tmux at 90x30 in the scratch clone. Picking the gate item opened the form: four ~300-char questions in a 5-row text block with 'pgup/pgdn or wheel scroll · 1-4/20' above four TextAreas; PageDown moved through 5-8, 9-12, 13-16, 17-20/20 to END4 with the TextAreas still on screen. Typed four answers and submitted: four Decisions written in the plan change set (question.answer --run), then run.resume failed (journal has no workflow definition) and the app showed 'answered; resume failed: ...' in a detail panel that scrolled with arrows and PageDown (1-10/25 to 16-25/25 at 90x15). Blank answers re-asked with 'answer each question'. With the worktree unreadable the form shows 'question 1 of 4: This question's text could not be read.' and no handle on screen. Arrow keys go to the TextArea and do not scroll the question text. Mouse wheel not exercised by hand (Go test only).

## Links

### Out

- verifies [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
