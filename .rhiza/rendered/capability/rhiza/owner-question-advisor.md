# owner-question-advisor

capability · `01M6PN1SZC653J5D55644NRYXE` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/run with a fake runner: leaving a plan-review, leaf-failed or cycle-cap question starts one advisor on that question and an advisor entry is journalled; a testscript (fake claude) run recommend on such a question records the restatement, recommendation and reasons, and rhiza inbox --json then shows them on the question's item; describe-golden regenerated for run.recommend.

## statement

When a run leaves a question for the owner from a Go template (a plan-review finding carried to the gate, a failed leaf, a review cycle cap), it starts the advisor on that question node as it does for a session's blocking question; the advisor records through rhiza run recommend a plain-language restatement of what is asked (no handles, IDs or finding IDs) with its recommendation and reasons, journalled against the question, while the template text stays the node's statement as the engineering record.

## Links

### Out

- decided-by [owner-question-wording-by-advisor](../../decision/owner-question-wording-by-advisor.md)

### In

- verifies [owner-question-advisor-verified](../../evidence/owner-question-advisor-verified.md)
- verifies [owner-question-advisor-qa-1](../../evidence/owner-question-advisor-qa-1.md)
- verifies [owner-question-advisor-verified-2](../../evidence/owner-question-advisor-verified-2.md)
- verifies [owner-question-advisor-qa-2](../../evidence/owner-question-advisor-qa-2.md)
- verifies [pr-313-fix-the-inbox-app-s-question](../../evidence/pr-313-fix-the-inbox-app-s-question.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
