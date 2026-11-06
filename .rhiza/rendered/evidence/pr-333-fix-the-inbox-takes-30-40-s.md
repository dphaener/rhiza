# pr-333-fix-the-inbox-takes-30-40-s

evidence · `01M6S6AMR6TXEF28BD315Q587P` · actor [dphaener](../actor/dphaener.md)

## head

22376617499cb57ea29801eb26b205306c4f1e49

## pr

https://github.com/dphaener/rhiza/pull/333

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #333 (https://github.com/dphaener/rhiza/pull/333) passed CI at 2237661

## status

verified

## Links

### Out

- verifies [inbox-gate-read-fast](../capability/rhiza/inbox-gate-read-fast.md)
- verifies [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)
- verifies [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
