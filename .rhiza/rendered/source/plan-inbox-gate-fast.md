# plan-inbox-gate-fast

source · `01M6RXX2NGGACDHQPP1TDAE3PW`

## quotes

```json
[]
```

## raw_sha256

af2dccc5ba5503240c22ff4c6e573d88871f54b1cfe8057170a25f4a013ad756

## session

fa508624-2462-49e9-98cd-01a1b1d9c450

## summary

Autonomous plan run for change set 01M6RXK2 (the inbox takes 30-40 s when a plan run is stopped at its gate; the Inbox app cannot scroll or act on that gate item). The planner read internal/api/inbox.go (inboxGateQuestions clears IndexPath and folds a model per question times per planned change set), question.go (gateAnswerCtx) and internal/tui/inboxapp (openGate/answerGate exist but are skipped when an item's Questions is empty, falling back to a non-scrollable detail panel). It raised four choices to the arbiter: the intent, the cheap read, the proof, and the app's action. The arbiter's Decision inbox-gate-questions-cheap settled them: serve milestone-26; read gate questions through the run worktree's own index with each planned change set folded at most once; always return one gate question per stop-detail handle even when the worktree's model can't be read; prove it with an internal/api budget test over the large generated log (no full-log decode, one fold per planned change set, within the read budget); in the Inbox app a gate item always opens the answer form, whose text scrolls, and every detail panel scrolls. No owner answered, so no owner words are quoted.

## title

Plan of the slow inbox at a gate, autonomous run with the arbiter

## Links

### In

- derived-from [inbox-gate-questions-cheap](../decision/inbox-gate-questions-cheap.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
