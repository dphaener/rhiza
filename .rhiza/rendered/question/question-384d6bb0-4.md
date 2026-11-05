# question-384d6bb0-4

question · `01M6J747QE3NSH03CFKXV23JC8`

## session

384d6bb0-11b5-42bf-a242-01a1a467f521

## statement

The change set's intent also mentions that each rhiza apply write takes ~6 s on the VM (31 writes > 3 min) and points to 'the daemon CPU bug'. No open change set for that bug exists, and the fix (write-path or daemon CPU performance) has nothing to do with interview capture or m2-interviews-recorded. I recommend keeping it out of this change set: park it as a question in this change set, linked to nothing it blocks, so the gate raises it and the owner opens its own change set for it. Planning it here would mix a performance investigation into a small capture fix.

## Links

### In

- answers [apply-latency-out-of-scope](../decision/apply-latency-out-of-scope.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
