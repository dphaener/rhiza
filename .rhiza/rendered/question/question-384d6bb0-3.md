# question-384d6bb0-3

question · `01M6J70FXQVNJRPEEEPF2JF6QS`

## session

384d6bb0-11b5-42bf-a242-01a1a467f521

## statement

Which side records an answer, and how is a proxy answer worded? Rhiza cannot tell who answered an AskUserQuestion in the pane (the session record has no answered-by field), so captured picks always get 'picked by the owner' wording. I recommend: the guide keeps the run's capture as the default recorder of picks (the planner writes no Decision for a plain owner pick); a planner writes the Decision itself only when the capture's wording would be wrong: the answer came from a proxy (an orchestrator or arbiter answering for the owner) or it carries a ruling the capture can't word, using the proxy's wording (e.g. the arbiter wording) and the verbatim question field, so the capture skips it. No new answered-by field in the session record now.

## Links

### In

- answers [capture-records-picks-planner-records-proxies](../decision/capture-records-picks-planner-records-proxies.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
