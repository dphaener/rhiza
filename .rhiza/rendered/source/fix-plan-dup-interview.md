# fix-plan-dup-interview

source · `01M6J7934ZA6AEM4FHWD9ZA2M6`

## quotes

```json
[]
```

## raw_sha256

958130c6b24413daade65a83cf1c7c40d372e75bba34b914f7fc10c330e484cd

## session

384d6bb0-11b5-42bf-a242-01a1a467f521

## summary

Autonomous plan run of change set 01M6J6FDD9GXDW4AENMS6YNY63; no owner present, every question went to the arbiter. Settled: the change set serves the landed done point m2-interviews-recorded; session interview/session done skip an answer when a Decision the change set created (or, once landed, a landed Decision derived-from the session's Source) has a question field equal to the asked question verbatim; the run's capture stays the default recorder of plain owner picks and the planner writes a Decision itself only for proxy answers or rulings the capture can't word, with proxy wording and the verbatim question; no answered-by session field now. The ~6 s per rhiza apply write on the VM is out of scope here and gets its own change set. Leaves planned: capture-skips-recorded-answers and plan-guide-who-records.

## title

Plan of 'Fix Plan interviews duplicate', autonomous with the arbiter

## Links

### In

- derived-from [fix-plan-dup-serves-interviews-recorded](../decision/fix-plan-dup-serves-interviews-recorded.md)
- derived-from [skip-recorded-answer-by-question](../decision/skip-recorded-answer-by-question.md)
- derived-from [capture-records-picks-planner-records-proxies](../decision/capture-records-picks-planner-records-proxies.md)
- derived-from [apply-latency-out-of-scope](../decision/apply-latency-out-of-scope.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
