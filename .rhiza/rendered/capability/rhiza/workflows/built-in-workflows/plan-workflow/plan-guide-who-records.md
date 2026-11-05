# plan-guide-who-records

capability · `01M6J7BMWSJXFQXYYDNRV7XD66` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A Go guide test finds in rhiza guide plan the rule that capture records plain owner picks with no planner Decision, the proxy exception with its wording, and the rule to copy the asked question verbatim into the question field; the guide text has no line telling the planner to write Decisions for plain picks

## statement

rhiza guide plan says the run's capture records the owner's AskUserQuestion picks and the planner writes no Decision for a plain owner pick; the planner writes one itself only when the capture's wording would be wrong (a proxy such as an orchestrator or arbiter answered, or the answer carries a ruling the capture cannot word), with the proxy's wording and the asked question verbatim in its question field so the capture skips that pick

## Links

### Out

- decided-by [capture-records-picks-planner-records-proxies](../../../../../decision/capture-records-picks-planner-records-proxies.md)
- decided-by [skip-recorded-answer-by-question](../../../../../decision/skip-recorded-answer-by-question.md)

### In

- verifies [plan-guide-who-records-verified](../../../../../evidence/plan-guide-who-records-verified.md)
- verifies [plan-guide-who-records-qa-1](../../../../../evidence/plan-guide-who-records-qa-1.md)
- verifies [plan-guide-who-records-qa-2](../../../../../evidence/plan-guide-who-records-qa-2.md)
- verifies [pr-281-fix-plan-interviews-duplicate](../../../../../evidence/pr-281-fix-plan-interviews-duplicate.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
