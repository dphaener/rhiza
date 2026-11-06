# inbox-gate-questions-two-cases

decision · `01M6RYN60X4SHND9CNN5WRWA16`

## reasoning

The review finding and critic agree the acceptance as written cannot be set up: in inboxGateQuestions (internal/api/inbox.go) a worktree that does not open returns before any handle is resolved, so 'worktree does not read' and 'a handle resolves in no planned change set' are mutually exclusive in one run, and 'three questions' vs 'both handles' does not add up. The leaf's statement names exactly these two fallbacks (unreadable model; unresolved handle), each a different code path, so each needs its own runnable case. Two questions per case keeps 'both handles' consistent and, in case (b), proves the unresolved fallback does not wipe the resolved question's statement. This is the split the finding proposed; it keeps the leaf's intent and adds no scope.

## statement

inbox-gate-questions-kept is accepted by two separate Go test cases in internal/api, not one combined setup: (a) a plan run stopped at its gate on two questions whose run worktree does not read gets an inbox open item listing both stop-detail handles, each with an empty statement, and its answer action; (b) a plan run stopped at its gate on two questions whose worktree reads, where one handle resolves in no planned change set, gets an inbox open item listing both handles: the unresolved one with an empty statement, the resolved one with its full statement. The existing testscript inbox-gate-questions still passes (go test -trimpath ./cmd/rhiza -run TestScripts/inbox-gate-questions).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6rxk0-r3-f1](../question/review-01m6rxk0-r3-f1.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
