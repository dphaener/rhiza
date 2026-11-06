# question-fa508624-1

question · `01M6RXQVR420SBRR5BA8PYARFG`

## session

fa508624-2462-49e9-98cd-01a1b1d9c450

## statement

Change set 01M6RXK2 (inbox takes 30-40 s at a gate; Inbox app can't act on the gate item) serves no intent yet. Which intent should it serve? / How should the inbox read a gate's open questions cheaply (today inboxGateQuestions opens the run worktree with no index, decoding the full 11,400-entry log, and folds a model per question × per planned change set)? / What proves the inbox stays fast at a gate? / In the Inbox app, what should a gate ('open') item do when picked?

## Links

### In

- answers [inbox-gate-questions-cheap](../decision/inbox-gate-questions-cheap.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
