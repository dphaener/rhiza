# fix-plan-dup-serves-interviews-recorded

decision · `01M6J6SHEQ39H4W35RB4HEB1RD`

## reasoning

m2-interviews-recorded promises that interview answers are recorded as Decisions automatically, with picked options as the agent's wording and only typed text quoted as the owner's. The change set fixes exactly two failures of that promise: session done writing duplicate answer-<session>-N Decisions for picks the planner already recorded, and owner wording applied to proxy answers. The code under repair (session done's answer capture, leaf done-captures-answers) is that done point's lineage, so no new intent is needed. The apply-latency note in the intention points at a separate daemon CPU bug and does not change which intent this fix serves.

## statement

Change set 'Fix Plan interviews duplicate' serves the landed done point m2-interviews-recorded.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-384d6bb0-1](../question/question-384d6bb0-1.md)
- derived-from [fix-plan-dup-interview](../source/fix-plan-dup-interview.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
