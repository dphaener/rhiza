# inbox-app-gate-scroll-qa-2

evidence · `01M6S5M9KCN02HMXVZA8YG1JWP` · actor [claude-code-9d02afc1](../actor/claude-code-9d02afc1.md)

## result

PASS

## statement

QA at 9efc1e89 with the built binary (go build -trimpath) in a scratch clone of this repo (11,630 log entries, worktree index present, load average 14 on 4 cores), a draft plan change set with four ~330-char open questions and a hand-written plan run journal stopped at its gate on them. rhiza question answer <q> <text> --run <run>: 1.3 s and 2.0 s with the new binary, 34 s with the binary at edf9ff99 (before the gateAnswerCtx fix). rhiza tui app inbox under tmux at 90x30: picking the gate item opened the form (pgup/pgdn or wheel scroll, 1-4/24; PageDown reached 17-20/24 with the four TextAreas on screen); typing four answers and submitting returned in 5.8 s (about 100 s in cycle 1) with four Decisions answer-qa3-q1..q4 in the plan change set, then 'answered; resume failed: ... workflow has no steps' (the hand-made journal has no workflow) in a detail panel that scrolled with PageDown (5-29/29). The answers written through the worktree index read the same with the log-only edf9ff99 binary (show qa3-q4: answered by answer-qa3-q4). go test -trimpath ./internal/api -run Gate and ./internal/tui/inboxapp/... pass.

## Links

### Out

- verifies [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
