# question-384d6bb0-2

question · `01M6J6W6Z5N8B4CH9QVPHPYGFV`

## session

384d6bb0-11b5-42bf-a242-01a1a467f521

## statement

How should session interview/session done tell that an AskUserQuestion answer is already recorded, so it skips it instead of writing a duplicate answer-<sid8>-N Decision? I recommend: skip an answer when a live Decision the change set created (or, when the change set has landed, a landed Decision derived-from this session's Source) has a `question` field equal to the asked question verbatim. It is exact, needs no new link type, and the planner can always copy the question text it asked. The alternative of skipping all capture whenever the planner wrote any Decision linked to the session's Source is coarser: picks the planner did not write would be lost.

## Links

### In

- answers [skip-recorded-answer-by-question](../decision/skip-recorded-answer-by-question.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
