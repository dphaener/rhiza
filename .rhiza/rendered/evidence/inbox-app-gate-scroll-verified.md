# inbox-app-gate-scroll-verified

evidence · `01M6S35P1XTAMF8ZYJ65JGWNRM` · actor [claude-code-52cc90d6](../actor/claude-code-52cc90d6.md)

## result

PASS

## statement

At 63456160 ran go test -trimpath -count=1 ./internal/tui/inboxapp/... (ok, fake runner only; vet and gofmt clean) and read the diff. TestGateFormScrollsLongQuestions: four ~400-char questions in a 70x30 window show the form within 30 rows with END4 hidden; page-down and then the mouse wheel (MouseMode on) reveal END4 while 'Your answer to question 1 of 4' stays visible; four answers then submit 1+4+1 calls ending in run.resume and quit. TestGateUnreadableQuestionShowsPositionNotHandle: an empty statement shows 'question 2 of 4: This question's text could not be read.' and 'Your answer to question 2 of 4', with no q-handle anywhere. TestGateAnswersEachThenResumes checks question.answer <handle> <text> --run R4 per question then run.resume; TestGateBlankAnswerSendsNothing: blank calls nothing. TestGateAnswerFailureReopensOnRemaining: a failed answer reopens on the remaining questions with 'answer failed: disk full' and no run.resume. TestGateResumeFailureDetailScrolls: a failed resume shows 'resume failed: busy' in a ui.InfoPanel sized to the window (wheel-scrolled by Model.Update). TestDetailPanelScrolls: a long stalled-item detail is clipped, page-down reaches LASTWORD, wheel-up scrolls back.

## Links

### Out

- verifies [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
