# inbox-gate-questions-kept

capability · `01M6RXXHCHHB08YWNFGY67DT82` · parent [rhiza](../rhiza.md)

## acceptance

Two Go tests in internal/api (decision inbox-gate-questions-two-cases): (a) a plan run stopped at its gate on two questions whose run worktree does not read gets an inbox open item listing both stop-detail handles, each with an empty statement; (b) a plan run stopped at its gate on two questions whose worktree reads, where one handle resolves in no planned change set, lists both handles, the unresolved one with an empty statement and the resolved one with its full statement. The existing testscript inbox-gate-questions still passes (go test -trimpath ./cmd/rhiza -run TestScripts/inbox-gate-questions).

## statement

A gate (open) inbox item always carries its questions: when the run worktree's model can't be read or a handle doesn't resolve, the inbox still lists one gate question per handle of the stop detail, with its statement empty, so the item never loses its questions or its answer action.

## Links

### Out

- decided-by [inbox-gate-questions-cheap](../../decision/inbox-gate-questions-cheap.md)

### In

- blocks [review-01m6rxk0-r3-f1](../../question/review-01m6rxk0-r3-f1.md)
- verifies [inbox-gate-questions-kept-verified](../../evidence/inbox-gate-questions-kept-verified.md)
- verifies [inbox-gate-questions-kept-qa-1](../../evidence/inbox-gate-questions-kept-qa-1.md)
- verifies [inbox-gate-questions-kept-verified-2](../../evidence/inbox-gate-questions-kept-verified-2.md)
- verifies [pr-333-fix-the-inbox-takes-30-40-s](../../evidence/pr-333-fix-the-inbox-takes-30-40-s.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
