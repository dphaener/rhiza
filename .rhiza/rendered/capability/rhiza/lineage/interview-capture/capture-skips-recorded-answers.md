# capture-skips-recorded-answers

capability · `01M6J7B10MH74CTH61H69ETSDW` · parent [interview-capture](../interview-capture.md)

## acceptance

A testscript whose fake agent leaves a transcript with two picked answers, in a change set that already holds a Decision whose question field is the first question verbatim: session done proposes only answer-<sid8>-2 and reports answer-<sid8>-1 as skipped; a Decision whose question differs by one character does not stop the proposal; and with the change set landed and a landed Decision derived-from the session's Source carrying the first question, the record capture writes only the second answer

## statement

rhiza session interview and session done skip an answered question when a live Decision the session's change set created (or, once that change set has landed, a live landed Decision derived-from the session's Source) has a question field equal to the asked question verbatim, and report it as skipped; only answers with no such Decision become answer-<sid8>-N Decisions

## Links

### Out

- decided-by [skip-recorded-answer-by-question](../../../../decision/skip-recorded-answer-by-question.md)
- decided-by [capture-records-picks-planner-records-proxies](../../../../decision/capture-records-picks-planner-records-proxies.md)

### In

- verifies [capture-skips-recorded-answers-verified](../../../../evidence/capture-skips-recorded-answers-verified.md)
- verifies [capture-skips-recorded-answers-qa-1](../../../../evidence/capture-skips-recorded-answers-qa-1.md)
- verifies [capture-skips-recorded-answers-verified-2](../../../../evidence/capture-skips-recorded-answers-verified-2.md)
- verifies [capture-skips-recorded-answers-qa-2](../../../../evidence/capture-skips-recorded-answers-qa-2.md)
- verifies [pr-281-fix-plan-interviews-duplicate](../../../../evidence/pr-281-fix-plan-interviews-duplicate.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
