# inbox-app-answer-shows-question-verified

evidence · `01M6PQT5YS201NJ959VRW51KNR` · actor [claude-code-b25460f6](../actor/claude-code-b25460f6.md)

## result

PASS

## statement

Verifier at e9e08fd3: ran go test -trimpath ./internal/tui/inboxapp/... -count=1 -v; all 31 tests PASS with the fake runner, including the new TestGateFormShowsEachQuestionNoHandleAndAnswersRightOne: a two-question gate form shows '1 of 2: <restatement>' and 'Recommended: Keep it' above 'Your answer to question 1 of 2', no handle (review-01m6pf30-r3-f1, claude-design-output) anywhere in either step's view, and submitting calls question.answer with [review-01m6pf30-r3-f1 yes --run R4] then [claude-design-output no --run R4]. For the single answer form, a throwaway (deleted, uncommitted) test with a 300+ char question item at width 70 showed every word wrapped, then 'Recommended:' and 'Why:', all above the 'Your answer' TextArea; TestBlockerShowsEveryWordWrapped and TestNoSurfaceShowsAnIDOrHandle also pass. Read the diff: test-only, matching askLines/openGate/answerGate in inboxapp.go.

## Links

### Out

- verifies [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
