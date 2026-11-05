# inbox-question-fix-plan

source · `01M6PN0RAQYP12ANPT84MWSMVN`

## quotes

```json
[
  {
    "by": "owner",
    "text": "The full text would have been more helpful than the internal handle. No need to ever show the internal handle, it's nonsense to a user.",
    "where": "message to the plan session, 2026-11-05"
  }
]
```

## raw_sha256

91c179fd238a5afa5422f89c28fab2d5aa79ab71bc5d666c3d52b077ccd56eaf

## session

4dd18cc3-6192-406a-b78d-01a1ad4cf530

## summary

Autonomous plan run on change set 'Fix The Inbox app's question'. Code reading found the gate form labels each answer box with the open question's handle (gateQuestions parses handles from the stop detail), blocker items name what they block by handle, the detail panel shows the run ID, and an overlap notice names the other change set by ID. The three Go-templated owner questions (plan review, leaf failed, review cycle cap) carry handles and finding IDs and no recommendation. The owner sent a message making it a rule that no owner-facing surface shows a handle or ID. The arbiter ruled that the run's existing advisor writes a plain-language restatement and recommendation for every templated question, and that the CLI's human text keeps handles and IDs but prints the full statement too. The change set serves the milestone-26 Inbox app intent.

## title

Plan of the Inbox question fix, autonomous run with the arbiter

## Links

### In

- derived-from [owner-surfaces-full-text](../decision/owner-surfaces-full-text.md)
- derived-from [owner-question-wording-by-advisor](../decision/owner-question-wording-by-advisor.md)
- derived-from [owner-surfaces-no-ids](../constraint/owner-surfaces-no-ids.md)

## History

- `01M6PMSN1PKES7GJYAEXNZSN13` Fix The Inbox app's question items aren't actionable: a row shows the question's handle or short title (e.g. 'review-01m6pf30-r3-f1', 'Claude Code Design') instead of what is asked, so the owner cannot tell what to answer. Seen 2026-11-05 on the web UI milestone's gate (run 01M6PF30, questions review-01m6pf30-r3-f1 and claude-design-output). Each question row should show the question's full statement (wrapped), its change set title, and what it blocks; the answer form should show the full statement above the text box. Also, questions raised for the owner (plan-review findings carried to the gate, leaf-failed questions) should be written for the owner in plain language with the agent's recommendation, not as an engineering finding: the review finding about Operation.Write had to be re-explained in chat before the owner could decide.
