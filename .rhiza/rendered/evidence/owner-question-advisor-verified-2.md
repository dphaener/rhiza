# owner-question-advisor-verified-2

evidence · `01M6PV1ANZBBDX2BJ9EG5HYZG5` · actor [claude-code-97098653](../actor/claude-code-97098653.md)

## result

PASS

## statement

Verifier, attempt 2 at 69988347. I read the fix diff: inboxBlockers now joins questionAdvice(runs), the latest advice with a recommendation per question node ID across runs, onto blocker items. These carry a new restatement field plus the recommendation and reasons. rhiza inbox prints them, and the Inbox app shows the restatement in the row and above the answer form. Ran: go test ./internal/run -run StartsAdvisor -v: TestLeafFailureStartsAdvisor, TestCycleCapQuestionStartsAdvisor and TestCommandAskedNormalRunStartsAdvisor (plan-review) pass. go test ./cmd/rhiza -run 'TestScripts/(inbox-.*|describe-golden|tui-inbox-app)$' is ok. That includes inbox-question-advisor (fake claude, run recommend, inbox --json shows the restatement, recommendation and reasons; the node's statement stays the template), the new inbox-blocker-advice (a leaf-failed question's blocker item shows restatement, recommendation and reasons in --json and text, with the template question kept) and describe-golden (run.recommend plus the new field). go test ./internal/tui/inboxapp (incl. TestBlockerShowsAdvisorRestatement) is ok. go test ./internal/api -run 'Inbox|Blocker|Advice|Recommend' is ok. An earlier unscoped internal/api run was killed by a signal in the timing-budget tests, not by a failed assertion. This fixes the attempt-1 finding f2: leaf-failed and cycle-cap advice now reaches the blocker item.

## Links

### Out

- verifies [owner-question-advisor](../capability/rhiza/owner-question-advisor.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
