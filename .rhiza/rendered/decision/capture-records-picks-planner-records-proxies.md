# capture-records-picks-planner-records-proxies

decision · `01M6J7218QNN4EP1B088PYQKTT`

## reasoning

The change set serves m2-interviews-recorded, which promises interview answers are recorded as Decisions automatically; keeping capture as the default preserves that promise and means a planner forgetting to write one loses nothing. The duplicates came from both sides recording the same picks, so the guide must name one default recorder, and capture is the one that already exists and words owner picks correctly. Capture's only wrong case is a proxy answer, which it cannot detect because the session record has no answered-by field; the side that knows (the planner, or the orchestrator acting as arbiter) writing the Decision with proxy wording and the verbatim question fixes that and reuses the question-field skip proposed in the sibling question, with no new schema. Adding an answered-by field would need hook/pane changes beyond this fix's scope and can come later if proxy answers become common.

## statement

The run's answer capture stays the default recorder of AskUserQuestion picks: the planner writes no Decision for a plain owner pick. A planner writes the Decision itself only when the capture's wording would be wrong (the answer came from a proxy such as an orchestrator or arbiter answering for the owner, or it carries a ruling the capture cannot word), using the proxy's wording (e.g. '(assumed, would have escalated)') and the asked question verbatim in its question field, so the capture skips that pick. No answered-by field is added to the session record now.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-384d6bb0-3](../question/question-384d6bb0-3.md)
- derived-from [fix-plan-dup-interview](../source/fix-plan-dup-interview.md)

### In

- decided-by [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)
- decided-by [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
