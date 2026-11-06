# review-01m6rxk0-r3-f1

question · `01M6RYKX4E14KYYST0X8X5AJMP`

## statement

Plan review round 3 left finding f1 (missing-acceptance on inbox-gate-questions-kept): The acceptance check mixes two cases into one setup and does not add up: 'a plan run stopped at its gate on two questions whose worktree does not read (and one whose handle resolves in no planned change set) gets an inbox open item listing both handles'. If the worktree does not read, no handle can resolve, so the second case cannot happen in the same run, and 'both' contradicts the three questions named. A builder cannot tell what the test should set up or assert. Split it into two runnable cases: (a) the worktree does not read, so every stop-detail handle is listed with an empty statement; (b) the worktree reads but one handle resolves in no planned change set, so that handle is listed with an empty statement and the others keep their statements. The critic confirmed it: inbox-gate-questions-kept's acceptance asks for one test where the run worktree does not read and also one handle resolves in no planned change set. In inboxGateQuestions (internal/api/inbox.go:1066-1071), a worktree that does not open returns before any handle is resolved, so both cases cannot happen in one run. It also names three questions but expects 'both handles' listed. As written, the setup can't be run as stated and leaves unclear which fallback path gets tested. Each of the leaf's two fallbacks needs its own runnable case. How should it be settled?

## Links

### Out

- blocks [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)

### In

- answers [inbox-gate-questions-two-cases](../decision/inbox-gate-questions-two-cases.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
