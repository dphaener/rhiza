# apply-latency-out-of-scope

decision · `01M6J760JG6SAFYCVBB7HZQD0G`

## reasoning

The change set serves m2-interviews-recorded (decision fix-plan-dup-serves-interviews-recorded) and its other decisions (skip-recorded-answer-by-question, capture-records-picks-planner-records-proxies) are all about answer capture in session done and the planner guide. Apply latency is a write-path or daemon CPU problem with different code, a different investigation and no link to interview capture; the intention text itself only points to it as 'see the daemon CPU bug'. No open change set covers it, but mixing a performance investigation into a small capture fix would widen review and delay the fix, so it belongs in its own change set rather than as leaves here.

## statement

The ~6 s per rhiza apply write on the VM is out of scope for change set 'Fix Plan interviews duplicate': no leaf here addresses apply latency or the daemon CPU bug. It stays recorded only as this answered question, and the owner opens its own change set for the write-path/daemon CPU performance work.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-384d6bb0-4](../question/question-384d6bb0-4.md)
- derived-from [fix-plan-dup-interview](../source/fix-plan-dup-interview.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
