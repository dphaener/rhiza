# skip-recorded-answer-by-question

decision · `01M6J6YC25BNHF4K7J7X3SHNG8`

## reasoning

The change set's intent asks session done to skip picks a Decision in the change set already answers. Decisions already have a question field (internal/op/fields.go) and session interview already holds each question verbatim (internal/api/session_interview.go), so an exact match needs no new link type or schema and is deterministic. Skipping all capture whenever the planner wrote any Source-linked Decision is coarser and would drop picks the planner did not record, losing owner answers, which is worse than an occasional duplicate. Covering the landed case through derived-from the session's Source keeps it correct when the change set lands early. Requiring planners to copy the question verbatim is a cheap guide rule that makes the match reliable.

## statement

session interview/session done skips an AskUserQuestion answer when a live Decision the change set created (or, once the change set has landed, a landed Decision derived-from this session's Source) has a question field equal to the asked question verbatim; only answers with no such Decision become answer-<sid8>-N Decisions. The planner guide tells planners who record an answer themselves to copy the asked question text verbatim into the Decision's question field.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-384d6bb0-2](../question/question-384d6bb0-2.md)
- derived-from [fix-plan-dup-interview](../source/fix-plan-dup-interview.md)

### In

- decided-by [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)
- decided-by [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
