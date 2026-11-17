# review-ruling-in-node-change-set

decision · `01M7MR431E9SSZAZ9BVDWC7T04`

## reasoning

The change set's intent is that a ruling is written and applied before merge, or the run fails; option (d) would leave every real autonomous plan run with a leaf ruling failed, which is the bug again in another form, and (c) would drop decision ruling-applied-means-decided-by, which this change set already built --rulings on. A link needs both nodes in one view, so the Decision must sit where the leaf is visible: the serving draft sees the intention's landed base but not its open ops, and the intention does not see the draft's ops. Option (a) alone fails too: the Decision's answers link would then not see a question left in the intention, so (a) needs (b). The code already has the placement rule for landed run change sets (reviewQuestionTarget picks the in-flight plan change set that created the node) and --rulings already searches in-flight plan views (reviewRulingsCheck), so extending the placement to open run change sets and carrying the target on the asked entry (run.Asked has no change field today; the arbiter gets st.Change in engine.go) is the smallest change that makes every view consistent. Additive journal field, so no breaking note.

## statement

A review question and its ruling live in the change set that holds the finding's node. rhiza run review-check puts the question on a finding whose node a plan change set still in flight created (a serving draft's leaf) into that change set, whether or not the run's change set has landed (reviewQuestionTarget no longer short-circuits to the run's change set while it is open), and journals that change set on the asked entry (an additive omitempty change field). The engine starts the arbiter for a session-less asked entry with RHIZA_CHANGE set to that change set, falling back to the run's change set when the entry names none, so the question, the arbiter's Decision and the leaf share one view and the reviser links the leaf decided-by the Decision there; --rulings already reads in-flight plan change sets' views. Findings on landed nodes or nodes of the run's own change set keep today's placement. Not option (c) or (d): no witness other than a decided-by link, and the script is not narrowed; plan-workflow-autonomous-left.txtar's leaf ruling on export-csv must apply and pass.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-013bd336-1](../question/question-013bd336-1.md)
- answers [ruling-on-serving-draft-node-unappliable](../question/ruling-on-serving-draft-node-unappliable.md)

### In

- decided-by [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)
- decided-by [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
