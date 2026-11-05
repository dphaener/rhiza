# m2-interviews-recorded

intent · `01M5SZHJQWABTT5TWCFMPTJDEG` · parent [milestone-2](../milestone-2.md)

## statement

Interview answers are recorded as Decisions automatically, a picked option as the agent's wording accepted by the owner and only typed text quoted as the owner's.

## Links

### In

- serves change set `01M6J6FDD9GXDW4AENMS6YNY63` (Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug)
- relates-to [m2-interviews-evidence](../../../evidence/m2-interviews-evidence.md) — meets this done point

## History

- `01M5SZHJ4AECQ3Z6PE60XFQZP6` Milestone 2 planning record: the conversation, the milestone-2 intent and its done points, and its decisions
