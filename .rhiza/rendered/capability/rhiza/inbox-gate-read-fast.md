# inbox-gate-read-fast

capability · `01M6RXXHAXYY47KE0DZ80E9B1X` · parent [rhiza](../rhiza.md)

## acceptance

A new budget test in internal/api over the large generated log with an index puts a plan run at its gate with several open questions spread over several planned change sets and asserts that inbox returns every question with its statement, decodes no full log (the worktree index catch-up decodes 0 entries when fresh), folds each planned change set at most once, and answers within the read budget; run with go test -trimpath ./internal/api -run <that test> plus the existing inbox tests in internal/api.

## statement

The inbox reads a plan run's gate questions cheaply: inboxGateQuestions opens the run worktree's model through its own index (index.DefaultPath of the worktree root) instead of decoding the full disk log, and folds each planned change set's model at most once per call, resolving every open handle against it; rhiza inbox, the daemon poll's gate-question fill, the Inbox app and the web UI's Overview and inbox fragment all go through this read and stay well under 1 s at a gate.

## Links

### Out

- decided-by [inbox-gate-questions-cheap](../../decision/inbox-gate-questions-cheap.md)

### In

- verifies [inbox-gate-read-fast-verified](../../evidence/inbox-gate-read-fast-verified.md)
- verifies [inbox-gate-read-fast-qa-1](../../evidence/inbox-gate-read-fast-qa-1.md)
- verifies [inbox-gate-read-fast-qa-2](../../evidence/inbox-gate-read-fast-qa-2.md)
- verifies [pr-333-fix-the-inbox-takes-30-40-s](../../evidence/pr-333-fix-the-inbox-takes-30-40-s.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
