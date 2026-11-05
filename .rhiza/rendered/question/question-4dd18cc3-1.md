# question-4dd18cc3-1

question · `01M6PMYQY84J91NMR5XNSHBJMY`

## session

4dd18cc3-6192-406a-b78d-01a1ad4cf530

## statement

Questions the runs raise for the owner (plan-review findings left at the gate, leaf-failed questions, review cycle-cap questions) are written today by Go templates from engineering findings, e.g. 'Plan review round 3 left finding f1 (gap on <handle>): ...'. To word them in plain language with a recommendation, who writes that wording? / The rule 'no owner-facing surface shows a node's handle or ID' covers the Rex apps (Inbox rows, answer form, detail panel), notifications and the web UI. Does it also cover the CLI's human text (e.g. rhiza inbox, rhiza run status)?

## Links

### In

- answers [owner-question-wording-by-advisor](../decision/owner-question-wording-by-advisor.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
