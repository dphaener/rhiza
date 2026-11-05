# owner-question-wording-by-advisor

decision · `01M6PN015M9RH75DV0K2A125CF`

## reasoning

Read the question, the change set intent and the code. internal/api/run_review.go reviewQuestionText builds the plan-review question from the finding with fmt.Sprintf; a Go template cannot write plain language or weigh a recommendation, so an agent must. The run engine already has exactly that agent: internal/run/engine.go advise() starts an advisor (rhiza guide advisor) that writes a recommendation with reasons via rhiza run recommend, and internal/api/inbox.go already carries Recommendation/Reasons on question items. Reusing it for template-raised questions adds no new session kind and gives the owner the recommendation the intent asks for, while keeping the engineering finding intact for agents and audit. Making the raising step's own agent word it was rejected: the gate questions are written by Go after the critic, not by an agent session. For the CLI: the owner's goal is to run Rhiza from the Rex apps (owner-goal-rex-only); the CLI is the agents' and developers' interface, where handles and IDs are how nodes are addressed, so removing them would break its use; the intent's failure (a row showing only a handle) is fixed there by also printing the full statement.

## statement

(1) The run's existing advisor session writes the owner-facing wording. For every question a run leaves for the owner from a Go template (plan-review findings carried to the gate, leaf-failed questions, review cycle-cap questions), the run starts the advisor on that question node, as it already does for a session's blocking question; the advisor reads the slice, the finding and the code and records, through rhiza run recommend, a plain-language restatement of what is being asked plus its recommendation and reasons. The Inbox row, answer form and notifications show that restatement and recommendation; the template text stays as the node's statement, the engineering record, viewable as detail. Until the advisor has written, surfaces show the template statement in full. (2) The no-handle-or-ID rule does not cover the CLI's human text: rhiza inbox, rhiza run status and other CLI output stay agent and developer surfaces and keep handles and IDs, since commands take them as arguments; they should also print the question's full statement and recommendation, not the handle alone.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-4dd18cc3-1](../question/question-4dd18cc3-1.md)
- derived-from [inbox-question-fix-plan](../source/inbox-question-fix-plan.md)

### In

- decided-by [inbox-question-text](../capability/rhiza/inbox-question-text.md)
- decided-by [owner-question-advisor](../capability/rhiza/owner-question-advisor.md)
- decided-by [inbox-app-answer-shows-question](../capability/rhiza/inbox-app-answer-shows-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
