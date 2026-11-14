# uncapped-loop-is-plain-red-loop

decision · `01M7DW59KBAFVF5QZQSTS1B26N`

## question

Change set 01M79KHB: rhiza workflow check should warn about an 'uncapped hand-back loop', but every handback:<step> route is already capped by handback_cap (at least 1, default 3). What should count as uncapped?

## reasoning

Every handback:<step> route is already capped by handback_cap; a plain red loop is the only loop it does not count, and it delivers no handback.md.

## statement

An uncapped hand-back loop is a red outcome (any on.<outcome> but pass) that names, by plain step name rather than handback:<agent step>, a step leading back to the routing step through an agent step: handback_cap does not count it, no handback.md is delivered and no inbox item parks it.

## wording

Plain red loop (assumed, would have escalated)

## Links

### Out

- derived-from [plan-skill-check](../source/plan-skill-check.md)

### In

- decided-by [workflow-check-uncapped-loop](../capability/rhiza/workflows/workflow-authoring/workflow-check-uncapped-loop.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
