# owner-question-advisor-verified

evidence · `01M6PRACVMPTMPKJYFYWHQGEG4` · actor [claude-code-654cd4ce](../actor/claude-code-654cd4ce.md)

## result

PASS

## statement

Verified at d0200787 (leaf commits 2c0166a7 and d0200787). I read the diff. In a normal run, adviseQuestion starts one advisor on a question node and journals an advisor entry whose Asker is the node ID. It is called from endLeaf (failed leaf), from capped (review cycle cap) and from arbitrateAsked (a session-less asked entry from a command step, as run review-check writes for plan-review questions; each question is advised once). run recommend takes --restatement, journalled as restated, and the inbox gate item shows it as restatement / 'in plain words', with the template statement unchanged. Ran go test -trimpath ./internal/run/... ./internal/api/...: ok. TestLeafFailureStartsAdvisor, TestCycleCapQuestionStartsAdvisor and TestCommandAskedNormalRunStartsAdvisor use a fake runner. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-question-advisor|inbox-gate-questions|inbox-advisor|describe-golden|run-review.*)$': all PASS. In inbox-question-advisor.txtar a fake claude advisor launched on a plan-review-style question node runs rhiza run recommend with restatement, recommendation and reasons. One advisor entry and the recommended entry are journalled against the node ID, and rhiza inbox --json shows restatement, recommendation and reasons on the gate-q item. describe-golden includes run.recommend's restatement. Both earlier findings (f1, f2) are addressed.

## Links

### Out

- verifies [owner-question-advisor](../capability/rhiza/owner-question-advisor.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
