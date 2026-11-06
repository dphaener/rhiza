# inbox-app-gate-scroll

capability · `01M6RXXHH0SGR3DGTJTD37G584` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/inboxapp with a fake runner: a gate item with four long questions in a short window shows the form, and page-down and the mouse wheel reveal the last question's text while the TextAreas stay reachable; a gate item whose questions have empty statements still opens the form, showing 'This question's text could not be read.' with its position ('question 2 of 4') above each TextArea, and no handle anywhere in the view; submitting calls question.answer --run per question then run.resume and quits; a blank answer calls nothing; a failed question.answer reopens the form on the remaining questions with 'answer failed: <error>' and calls no run.resume; a failed resume shows a detail panel that scrolls; a long detail of another item kind scrolls (go test -trimpath ./internal/tui/inboxapp/...).

## statement

In the Inbox app, picking a gate (open) item always opens the gate answer form: the full text of every open question with the advisor's recommendation scrolls (arrow and page keys and the mouse wheel) within the modal's height above one TextArea per question; submitting runs question.answer <handle> <text> --run <run> for each, then run.resume, then closes; a blank answer asks again; a failed answer reopens the form on the questions not yet answered with 'answer failed: <error>' so the owner can retry, as answerGate does today, and only a failed resume shows a scrollable detail with the error; a question whose statement is empty (inbox-gate-questions-kept's fallback) is shown above its TextArea as 'This question's text could not be read.' with its position among the gate's questions ('question 2 of 4') and never its handle, so inbox-app-answer-shows-question and owner-surfaces-no-ids still hold; and the detail panel of every other item kind scrolls too.

## Links

### Out

- decided-by [inbox-gate-questions-cheap](../../decision/inbox-gate-questions-cheap.md)

### In

- verifies [inbox-app-gate-scroll-verified](../../evidence/inbox-app-gate-scroll-verified.md)
- verifies [inbox-app-gate-scroll-qa-1](../../evidence/inbox-app-gate-scroll-qa-1.md)
- verifies [inbox-app-gate-scroll-verified-2](../../evidence/inbox-app-gate-scroll-verified-2.md)
- verifies [inbox-app-gate-scroll-qa-2](../../evidence/inbox-app-gate-scroll-qa-2.md)
- verifies [pr-333-fix-the-inbox-takes-30-40-s](../../evidence/pr-333-fix-the-inbox-takes-30-40-s.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
