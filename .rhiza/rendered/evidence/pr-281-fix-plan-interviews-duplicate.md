# pr-281-fix-plan-interviews-duplicate

evidence · `01M6MXMF6Z4D1XEC0XXA6RM2GA` · actor [dphaener](../actor/dphaener.md)

## head

fa7253b9bf18ddc80673224be8bb05162e65e9af

## pr

https://github.com/dphaener/rhiza/pull/281

## result

1 check passed: check SUCCESS

## statement

Pull request #281 (https://github.com/dphaener/rhiza/pull/281) passed CI at fa7253b

## status

verified

## Links

### Out

- verifies [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)
- verifies [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
