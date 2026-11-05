# inbox-app-answer-shows-question-verified-2

evidence · `01M6PV8ND7JRH097RWYA2JG19Q` · actor [claude-code-62f69d10](../actor/claude-code-62f69d10.md)

## result

PASS

## statement

Verifier claude-code-62f69d10 at 059598cd (attempt 2): ran go test -trimpath -count=1 ./internal/tui/inboxapp/... (ok) plus -v on the answer/gate tests, all PASS, and read the diff. TestQuestionFormShowsStatementAndRecommendationAboveAnswer (new, fixing finding f1) asserts, for both question and blocker items over a fake runner, every word of the full statement, then 'Recommended: ...', then 'Your answer' in that order. TestGateFormShowsEachQuestionNoHandleAndAnswersRightOne shows each of two gate questions' text (and recommendation) above its own TextArea with neither handle in the view, and checks question.answer was called with each handle and its typed text (--run R4). TestAnswerFormOpenedAfterResizeWraps checks a form opened after a resize wraps within the pane. go vet and gofmt clean.

## Links

### Out

- verifies [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
