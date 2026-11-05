# question-384d6bb0-1

question · `01M6J6Q2REYYQEFTH48R7FCN3H`

## session

384d6bb0-11b5-42bf-a242-01a1a467f521

## statement

Change set 'Fix Plan interviews duplicate' serves no intent yet. Which intent should it serve? I recommend the landed done point m2-interviews-recorded ('Interview answers are recorded as Decisions automatically, a picked option as the agent's wording accepted by the owner…'): the duplicate answer-<sid>-N Decisions and the wrong owner wording for proxy answers are failures of exactly that promise, and its lineage leaf done-captures-answers is the code being fixed.

## Links

### In

- answers [fix-plan-dup-serves-interviews-recorded](../decision/fix-plan-dup-serves-interviews-recorded.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
